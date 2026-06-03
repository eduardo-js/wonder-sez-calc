package server_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/wonderlic-calc/backend/internal/config"
	"github.com/wonderlic-calc/backend/internal/httpx"
	"github.com/wonderlic-calc/backend/internal/server"
)

func newTestRouter(origins []string) *gin.Engine {
	cfg := config.Config{
		Addr:           ":8080",
		AllowedOrigins: origins,
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return server.NewRouter(cfg, logger)
}

func TestHealthProbes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestRouter([]string{"http://localhost:5173"})

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/healthz", http.StatusOK, `"status":"ok"`},
		{"/readyz", http.StatusOK, `"status":"ready"`},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			engine.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantStatus)
			}
			if !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Errorf("body %q does not contain %q", w.Body.String(), tc.wantBody)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestRouter([]string{"http://localhost:5173"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"not_found"`) {
		t.Errorf("body %q missing not_found code", w.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestRouter([]string{"http://localhost:5173"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/healthz", nil)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: got %d, want 405", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"code":"method_not_allowed"`) {
		t.Errorf("body %q missing method_not_allowed code", w.Body.String())
	}
}

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	allowedOrigin := "http://localhost:5173"
	disallowedOrigin := "http://evil.example.com"

	engine := newTestRouter([]string{allowedOrigin})

	t.Run("allowed origin preflight", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
		req.Header.Set("Origin", allowedOrigin)
		req.Header.Set("Access-Control-Request-Method", "GET")
		engine.ServeHTTP(w, req)

		got := w.Header().Get("Access-Control-Allow-Origin")
		if got != allowedOrigin {
			t.Errorf("ACAO header: got %q, want %q", got, allowedOrigin)
		}
	})

	t.Run("disallowed origin preflight", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/healthz", nil)
		req.Header.Set("Origin", disallowedOrigin)
		req.Header.Set("Access-Control-Request-Method", "GET")
		engine.ServeHTTP(w, req)

		got := w.Header().Get("Access-Control-Allow-Origin")
		if got == disallowedOrigin {
			t.Errorf("ACAO header must not be set to disallowed origin %q", disallowedOrigin)
		}
		// gin-contrib/cors returns 403 for disallowed preflight
		if w.Code != http.StatusForbidden {
			t.Errorf("status: got %d, want 403 for disallowed origin", w.Code)
		}
	})
}

func TestValidationIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := newTestRouter([]string{"http://localhost:5173"})

	type probeDTO struct {
		Name string `json:"name" binding:"required"`
	}

	engine.POST("/api/v1/_probe", func(c *gin.Context) {
		var dto probeDTO
		if !httpx.BindJSON(c, &dto) {
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": dto.Name})
	})

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "valid body",
			body:       `{"name":"hello"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing required field",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeValidation,
		},
		{
			name:       "malformed json",
			body:       `{bad`,
			wantStatus: http.StatusBadRequest,
			wantCode:   httpx.CodeBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/_probe", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantStatus)
			}
			if tc.wantCode != "" && !strings.Contains(w.Body.String(), tc.wantCode) {
				t.Errorf("body %q does not contain code %q", w.Body.String(), tc.wantCode)
			}
		})
	}
}
