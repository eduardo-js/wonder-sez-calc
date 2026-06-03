// Package middleware provides reusable Gin middleware: panic recovery, request
// logging, and CORS.
package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wonderlic-calc/backend/internal/httpx"
)

// Recovery returns middleware that recovers from panics in downstream handlers,
// logs the panic via logger, and writes a 500 internal error envelope.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "error", fmt.Sprintf("%v", r))
				httpx.WriteError(c, http.StatusInternalServerError, httpx.CodeInternal, "internal server error", nil)
				return
			}
		}()
		c.Next()
	}
}
