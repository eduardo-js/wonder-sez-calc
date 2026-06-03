package health_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wonderlic-calc/backend/internal/health"
)

func TestHealthHandlers(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	_, mux := health.NewHandler(nil, logger)

	tests := []struct {
		name            string
		method          string
		path            string
		wantStatus      int
		wantContentType string
		wantBody        map[string]string
	}{
		{
			name:            "healthz returns 200 ok",
			method:          http.MethodGet,
			path:            "/healthz",
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			wantBody:        map[string]string{"status": "ok"},
		},
		{
			name:            "readyz returns 200 ready",
			method:          http.MethodGet,
			path:            "/readyz",
			wantStatus:      http.StatusOK,
			wantContentType: "application/json",
			wantBody:        map[string]string{"status": "ready"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			// Status code.
			if res.StatusCode != tc.wantStatus {
				t.Errorf("status: got %d, want %d", res.StatusCode, tc.wantStatus)
			}

			// Content-Type header (prefix match to tolerate "; charset=utf-8").
			ct := res.Header.Get("Content-Type")
			if len(ct) < len(tc.wantContentType) || ct[:len(tc.wantContentType)] != tc.wantContentType {
				t.Errorf("Content-Type: got %q, want prefix %q", ct, tc.wantContentType)
			}

			// JSON body.
			var got map[string]string
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				t.Fatalf("decode body: %v", err)
			}

			for key, wantVal := range tc.wantBody {
				if gotVal, ok := got[key]; !ok {
					t.Errorf("body missing key %q", key)
				} else if gotVal != wantVal {
					t.Errorf("body[%q]: got %q, want %q", key, gotVal, wantVal)
				}
			}
		})
	}
}
