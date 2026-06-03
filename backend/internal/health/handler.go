// Package health provides HTTP handlers for liveness and readiness probes.
package health

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// healthResponse is the JSON body returned by the liveness probe.
type healthResponse struct {
	Status string `json:"status"`
}

// readyResponse is the JSON body returned by the readiness probe.
type readyResponse struct {
	Status string `json:"status"`
}

// Handler holds dependencies for the health and readiness HTTP handlers.
type Handler struct {
	logger *slog.Logger
}

// NewHandler constructs a Handler and registers /healthz and /readyz on mux.
// If mux is nil, a new http.ServeMux is created and returned.
func NewHandler(mux *http.ServeMux, logger *slog.Logger) (*Handler, *http.ServeMux) {
	if mux == nil {
		mux = http.NewServeMux()
	}

	h := &Handler{logger: logger}
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /readyz", h.Readyz)

	return h, mux
}

// Healthz handles GET /healthz — liveness probe.
// Responds 200 application/json {"status":"ok"}.
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, r, http.StatusOK, healthResponse{Status: "ok"})
}

// Readyz handles GET /readyz — readiness probe.
// Responds 200 application/json {"status":"ready"}.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, r, http.StatusOK, readyResponse{Status: "ready"})
}

// writeJSON encodes v as JSON and writes it to w with the given status code.
// On encoding failure the error is logged and a 500 is returned.
func (h *Handler) writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.logger.Error("failed to encode JSON response",
			"error", err,
			"path", r.URL.Path,
		)
	}
}
