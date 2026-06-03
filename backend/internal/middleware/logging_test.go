package middleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wonderlic-calc/backend/internal/middleware"
)

func TestRequestLogger(t *testing.T) {
	tests := []struct {
		name       string
		handler    gin.HandlerFunc
		method     string
		path       string
		wantStatus int
		wantLevel  string
	}{
		{
			name: "200 logs at INFO",
			handler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
			wantLevel:  "INFO",
		},
		{
			name: "500 logs at ERROR",
			handler: func(c *gin.Context) {
				c.Status(http.StatusInternalServerError)
			},
			method:     http.MethodPost,
			path:       "/boom",
			wantStatus: http.StatusInternalServerError,
			wantLevel:  "ERROR",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			r := gin.New()
			r.Use(middleware.RequestLogger(logger))
			r.Handle(tc.method, tc.path, tc.handler)

			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			var logLine map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &logLine), "unmarshal log (raw: %s)", buf.String())

			assert.Equal(t, tc.method, logLine["method"])
			assert.Equal(t, tc.path, logLine["path"])
			// status is emitted as float64 in JSON.
			status, ok := logLine["status"].(float64)
			require.True(t, ok, "status key not a number: %v", logLine["status"])
			assert.Equal(t, tc.wantStatus, int(status))
			assert.NotNil(t, logLine["latency"], "latency key missing from log")
			assert.Equal(t, tc.wantLevel, logLine["level"])
		})
	}
}
