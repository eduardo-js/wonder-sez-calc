package middleware_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wonderlic-calc/backend/internal/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRecovery(t *testing.T) {
	tests := []struct {
		name            string
		handler         gin.HandlerFunc
		wantStatus      int
		wantCode        string
		wantPassthrough bool
	}{
		{
			name: "panic produces 500 envelope",
			handler: func(c *gin.Context) {
				panic("something went wrong")
			},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "internal",
		},
		{
			name: "normal handler passes through",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			},
			wantStatus:      http.StatusOK,
			wantPassthrough: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var logBuf strings.Builder
			logger := slog.New(slog.NewJSONHandler(&logBuf, nil))

			r := gin.New()
			r.Use(middleware.Recovery(logger))
			r.GET("/test", tc.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			assert.True(t, strings.HasPrefix(w.Header().Get("Content-Type"), "application/json"),
				"Content-Type should have application/json prefix, got %q", w.Header().Get("Content-Type"))

			if !tc.wantPassthrough {
				var body map[string]any
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body), "unmarshal body")

				errObj, ok := body["error"].(map[string]any)
				require.True(t, ok, "body missing 'error' key")

				assert.Equal(t, tc.wantCode, errObj["code"])
				assert.NotEmpty(t, errObj["message"], "message should not be empty")
				assert.Contains(t, logBuf.String(), "panic recovered")
			}
		})
	}
}

// Ensure Recovery works with a discard logger (no-op path).
func TestRecovery_discardLogger(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.GET("/", func(c *gin.Context) { panic("oops") })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
