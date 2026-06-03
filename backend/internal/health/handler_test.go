package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wonderlic-calc/backend/internal/health"
)

func TestHealthHandlers(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	engine := gin.New()
	health.Register(engine)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   map[string]string
	}{
		{
			name:       "healthz returns 200 ok",
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"status": "ok"},
		},
		{
			name:       "readyz returns 200 ready",
			path:       "/readyz",
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"status": "ready"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			engine.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.wantStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), "application/json")

			var got map[string]string
			require.NoError(t, json.NewDecoder(res.Body).Decode(&got))

			for k, want := range tc.wantBody {
				assert.Equal(t, want, got[k], "body[%q]", k)
			}
		})
	}
}
