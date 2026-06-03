// Package httpx provides HTTP error envelopes and JSON binding helpers for gin handlers.
package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error codes (machine-readable).
const (
	CodeValidation       = "validation_failed"
	CodeBadRequest       = "bad_request"
	CodeNotFound         = "not_found"
	CodeMethodNotAllowed = "method_not_allowed"
	CodeInternal         = "internal"
	CodeCalculation      = "calculation_error"
)

// Error is the structured error payload returned in the "error" envelope field.
type Error struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// envelope wraps Error for the wire shape {"error": {...}}.
type envelope struct {
	Error Error `json:"error"`
}

// WriteError writes {"error": Error} as JSON with the given status and aborts the context.
func WriteError(c *gin.Context, status int, code, message string, fields map[string]string) {
	// Normalise empty map to nil so omitempty suppresses the field.
	if len(fields) == 0 {
		fields = nil
	}
	c.AbortWithStatusJSON(status, envelope{
		Error: Error{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	})
}

// NotFound writes a 404 not_found error envelope.
func NotFound(c *gin.Context) {
	WriteError(c, http.StatusNotFound, CodeNotFound, "not found", nil)
}

// MethodNotAllowed writes a 405 method_not_allowed error envelope.
func MethodNotAllowed(c *gin.Context) {
	WriteError(c, http.StatusMethodNotAllowed, CodeMethodNotAllowed, "method not allowed", nil)
}
