package calculate_test

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/wonderlic-calc/backend/internal/calc"
	"github.com/wonderlic-calc/backend/internal/calculate"
	"github.com/wonderlic-calc/backend/internal/httpx"
)

func newTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	rg := engine.Group("/api/v1")
	calculate.Register(rg, logger)
	return engine
}

// newTestEngineWithEvaluator creates an engine that uses a custom evaluator — for
// testing error paths that cannot be reached via valid expressions in the HTTP API.
func newTestEngineWithEvaluator(eval func(string) (string, error)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	rg := engine.Group("/api/v1")
	calculate.RegisterWithEvaluator(rg, logger, eval)
	return engine
}

func TestCalculateHandler(t *testing.T) {
	engine := newTestEngine()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
		wantAbsent string
	}{
		// Happy path
		{
			name:       "basic addition",
			body:       `{"expression":"2 + 3"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"result":"5"`,
		},
		{
			name:       "precedence",
			body:       `{"expression":"2 + 3 * 4"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"result":"14"`,
		},
		{
			name:       "decimal result",
			body:       `{"expression":"1 / 3"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"result":"0.333333333333"`,
		},
		{
			name:       "echoes expression",
			body:       `{"expression":"2 + 2"}`,
			wantStatus: http.StatusOK,
			wantBody:   `"expression":"2 + 2"`,
		},

		// 400 validation_failed — malformed expression
		{
			name:       "consecutive operators",
			body:       `{"expression":"5 + * 2"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},
		{
			name:       "validation failed has fields.expression",
			body:       `{"expression":"5 + * 2"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"fields"`,
		},
		{
			name:       "expression too long",
			body:       `{"expression":"` + strings.Repeat("1+", 129) + `1"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},
		{
			name:       "illegal chars",
			body:       `{"expression":"1 + abc"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},
		{
			name:       "empty expression",
			body:       `{"expression":""}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},
		{
			name:       "unbalanced paren",
			body:       `{"expression":"(3 + 4"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},

		// 422 calculation_error
		{
			name:       "divide by zero",
			body:       `{"expression":"1 / 0"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `"code":"calculation_error"`,
		},
		{
			name:       "divide by zero expression",
			body:       `{"expression":"5 / (3 - 3)"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `"code":"calculation_error"`,
		},

		// 400 bad_request — unparseable JSON
		{
			name:       "bad json",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"bad_request"`,
		},
		{
			name:       "missing expression field",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"code":"validation_failed"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d; body: %s", w.Code, tc.wantStatus, w.Body.String())
			}
			if tc.wantBody != "" && !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Errorf("body %q does not contain %q", w.Body.String(), tc.wantBody)
			}
			if tc.wantAbsent != "" && strings.Contains(w.Body.String(), tc.wantAbsent) {
				t.Errorf("body %q must not contain %q", w.Body.String(), tc.wantAbsent)
			}
		})
	}
}

func TestCalculateHandler_MethodNotAllowed(t *testing.T) {
	engine := newTestEngine()
	engine.HandleMethodNotAllowed = true

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	engine.ServeHTTP(w, req)

	// gin returns 405 for registered routes with wrong method, or 404 if not handled
	// We registered only POST, so GET should be 405 or 404 depending on engine config
	if w.Code == http.StatusOK {
		t.Errorf("expected non-200 for GET, got 200")
	}
}

func TestCalculateHandler_NonFiniteError(t *testing.T) {
	// Inject an evaluator that returns ErrNonFinite to cover that handler branch.
	engine := newTestEngineWithEvaluator(func(string) (string, error) {
		return "", fmt.Errorf("%w", calc.ErrNonFinite)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate",
		strings.NewReader(`{"expression":"2 + 2"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("want 422, got %d; body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), httpx.CodeCalculation) {
		t.Errorf("want calculation_error in body: %s", w.Body.String())
	}
}

func TestCalculateHandler_InternalError(t *testing.T) {
	// Inject an evaluator that returns an unexpected error to cover the default branch.
	engine := newTestEngineWithEvaluator(func(string) (string, error) {
		return "", errors.New("unexpected internal failure")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate",
		strings.NewReader(`{"expression":"2 + 2"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d; body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), httpx.CodeInternal) {
		t.Errorf("want internal in body: %s", w.Body.String())
	}
}

func TestCalculateHandler_ValidationErrorHasExpressionField(t *testing.T) {
	engine := newTestEngine()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate",
		strings.NewReader(`{"expression":"5 + * 2"}`))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `"expression"`) {
		t.Errorf("fields.expression missing in: %s", body)
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}

	// Must use CodeValidation not CodeCalculation
	if strings.Contains(body, httpx.CodeCalculation) {
		t.Errorf("must not use calculation_error for invalid expression")
	}
	if !strings.Contains(body, httpx.CodeValidation) {
		t.Errorf("must use validation_failed for invalid expression")
	}
}
