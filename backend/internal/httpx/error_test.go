package httpx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wonderlic-calc/backend/internal/httpx"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

func TestConvenienceHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fn         func(*gin.Context)
		wantStatus int
		wantCode   string
	}{
		{"NotFound", httpx.NotFound, http.StatusNotFound, httpx.CodeNotFound},
		{"MethodNotAllowed", httpx.MethodNotAllowed, http.StatusMethodNotAllowed, httpx.CodeMethodNotAllowed},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, w := newContext(http.MethodGet, "/test")
			tc.fn(c)

			if w.Code != tc.wantStatus {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantStatus)
			}
			if !c.IsAborted() {
				t.Error("expected context aborted")
			}
			var body struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Error.Code != tc.wantCode {
				t.Errorf("code: got %q, want %q", body.Error.Code, tc.wantCode)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		code       string
		message    string
		fields     map[string]string
		wantFields map[string]string
	}{
		{
			name:       "validation_failed with fields",
			status:     http.StatusBadRequest,
			code:       httpx.CodeValidation,
			message:    "validation failed",
			fields:     map[string]string{"name": "required"},
			wantFields: map[string]string{"name": "required"},
		},
		{
			name:    "bad_request nil fields omitted",
			status:  http.StatusBadRequest,
			code:    httpx.CodeBadRequest,
			message: "bad request",
			fields:  nil,
		},
		{
			name:    "not_found",
			status:  http.StatusNotFound,
			code:    httpx.CodeNotFound,
			message: "not found",
			fields:  nil,
		},
		{
			name:    "method_not_allowed",
			status:  http.StatusMethodNotAllowed,
			code:    httpx.CodeMethodNotAllowed,
			message: "method not allowed",
			fields:  nil,
		},
		{
			name:    "internal",
			status:  http.StatusInternalServerError,
			code:    httpx.CodeInternal,
			message: "internal server error",
			fields:  nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, w := newContext(http.MethodGet, "/test")

			httpx.WriteError(c, tc.status, tc.code, tc.message, tc.fields)

			if w.Code != tc.status {
				t.Errorf("status: got %d, want %d", w.Code, tc.status)
			}

			ct := w.Header().Get("Content-Type")
			const wantCT = "application/json"
			if len(ct) < len(wantCT) || ct[:len(wantCT)] != wantCT {
				t.Errorf("Content-Type: got %q, want prefix %q", ct, wantCT)
			}

			if !c.IsAborted() {
				t.Error("expected context to be aborted")
			}

			var body struct {
				Error struct {
					Code    string            `json:"code"`
					Message string            `json:"message"`
					Fields  map[string]string `json:"fields"`
				} `json:"error"`
			}
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			if body.Error.Code != tc.code {
				t.Errorf("code: got %q, want %q", body.Error.Code, tc.code)
			}
			if body.Error.Message != tc.message {
				t.Errorf("message: got %q, want %q", body.Error.Message, tc.message)
			}

			if tc.wantFields != nil {
				for k, v := range tc.wantFields {
					if got, ok := body.Error.Fields[k]; !ok {
						t.Errorf("fields missing key %q", k)
					} else if got != v {
						t.Errorf("fields[%q]: got %q, want %q", k, got, v)
					}
				}
			} else {
				if body.Error.Fields != nil {
					t.Errorf("fields: expected nil/omitted, got %v", body.Error.Fields)
				}
			}
		})
	}
}
