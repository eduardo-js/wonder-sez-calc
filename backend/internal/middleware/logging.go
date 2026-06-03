package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns middleware that logs each request as structured JSON:
// method, path, status, latency (and any c.Errors).
func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)

		args := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency.String(),
		}

		if status >= http.StatusInternalServerError {
			logger.Error("request", args...)
		} else {
			logger.Info("request", args...)
		}
	}
}
