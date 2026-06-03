// Package server wires the application's gin.Engine with middleware and routes.
package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wonderlic-calc/backend/internal/calculate"
	"github.com/wonderlic-calc/backend/internal/config"
	"github.com/wonderlic-calc/backend/internal/health"
	"github.com/wonderlic-calc/backend/internal/httpx"
	"github.com/wonderlic-calc/backend/internal/middleware"
)

// NewRouter builds the application's gin.Engine with the full middleware chain
// and all routes registered centrally.
func NewRouter(cfg config.Config, logger *slog.Logger) *gin.Engine {
	engine := gin.New()
	engine.HandleMethodNotAllowed = true

	engine.Use(
		middleware.Recovery(logger),
		middleware.RequestLogger(logger),
		middleware.CORS(cfg.AllowedOrigins),
	)

	health.Register(engine)

	// Versioned API group.
	apiV1 := engine.Group("/api/v1")
	calculate.Register(apiV1, logger)

	engine.NoRoute(func(c *gin.Context) {
		httpx.WriteError(c, http.StatusNotFound, httpx.CodeNotFound, "resource not found", nil)
	})

	engine.NoMethod(func(c *gin.Context) {
		httpx.WriteError(c, http.StatusMethodNotAllowed, httpx.CodeMethodNotAllowed, "method not allowed", nil)
	})

	return engine
}
