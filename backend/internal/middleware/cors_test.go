package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/wonderlic-calc/backend/internal/middleware"
)

func TestCORS(t *testing.T) {
	allowedOrigin := "http://localhost:5173"
	disallowedOrigin := "http://evil.example"

	newRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(middleware.CORS([]string{allowedOrigin}))
		r.POST("/api", func(c *gin.Context) { c.Status(http.StatusOK) })
		return r
	}

	tests := []struct {
		name            string
		method          string
		origin          string
		wantAllowOrigin string
		wantStatus      int
		checkMethods    bool
	}{
		{
			name:            "preflight allowed origin",
			method:          http.MethodOptions,
			origin:          allowedOrigin,
			wantAllowOrigin: allowedOrigin,
			wantStatus:      http.StatusNoContent,
			checkMethods:    true,
		},
		{
			name:            "preflight disallowed origin",
			method:          http.MethodOptions,
			origin:          disallowedOrigin,
			wantAllowOrigin: "",
			wantStatus:      http.StatusForbidden,
		},
		{
			name:            "simple request allowed origin",
			method:          http.MethodPost,
			origin:          allowedOrigin,
			wantAllowOrigin: allowedOrigin,
			wantStatus:      http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newRouter()
			req := httptest.NewRequest(tc.method, "/api", nil)
			req.Header.Set("Origin", tc.origin)
			if tc.method == http.MethodOptions {
				req.Header.Set("Access-Control-Request-Method", "POST")
				req.Header.Set("Access-Control-Request-Headers", "Content-Type")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			assert.Equal(t, tc.wantAllowOrigin, w.Header().Get("Access-Control-Allow-Origin"))

			if tc.checkMethods {
				allowMethods := w.Header().Get("Access-Control-Allow-Methods")
				assert.Contains(t, allowMethods, "POST", "Access-Control-Allow-Methods missing POST")
				assert.Contains(t, allowMethods, "GET", "Access-Control-Allow-Methods missing GET")
			}
		})
	}
}
