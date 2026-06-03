// Package health provides Gin handlers for liveness and readiness probes.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthResponse is the JSON body returned by the liveness probe.
type healthResponse struct {
	Status string `json:"status"`
}

// readyResponse is the JSON body returned by the readiness probe.
type readyResponse struct {
	Status string `json:"status"`
}

// Healthz is the liveness probe: 200 {"status":"ok"}.
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{Status: "ok"})
}

// Readyz is the readiness probe: 200 {"status":"ready"}.
func Readyz(c *gin.Context) {
	c.JSON(http.StatusOK, readyResponse{Status: "ready"})
}

// Register mounts the probes on the given router.
//
//	GET /healthz, GET /readyz
func Register(r gin.IRouter) {
	r.GET("/healthz", Healthz)
	r.GET("/readyz", Readyz)
}
