// Package calculate provides the HTTP handler for POST /api/v1/calculate.
package calculate

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/wonderlic-calc/backend/internal/calc"
	"github.com/wonderlic-calc/backend/internal/httpx"
)

// allowedCharsRE matches the full set of characters permitted in an expression.
// Allowed: digits, '.', '+', '-', '*', '/', '(', ')', spaces, and 'e'/'E' for
// scientific-notation literals (e.g. 1e40). The evaluator does the structural
// validation; this gate only rejects clearly out-of-charset input.
var allowedCharsRE = regexp.MustCompile(`^[0-9eE+\-*/().\ ]+$`)

// calculationRequest is the JSON body accepted by the calculate endpoint.
type calculationRequest struct {
	Expression string `json:"expression" binding:"required"`
}

// calculationResult is the JSON body returned on a successful evaluation.
type calculationResult struct {
	Result     string `json:"result"`
	Expression string `json:"expression"`
}

// evaluatorFunc is the signature for expression evaluation; injectable for tests.
type evaluatorFunc func(expr string) (string, error)

// Handler holds handler dependencies for the calculate endpoint.
type Handler struct {
	logger    *slog.Logger
	evaluator evaluatorFunc
}

// newHandler constructs a Handler using the default calc.Evaluate implementation.
func newHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger, evaluator: calc.Evaluate}
}

// Calculate handles POST /calculate.
func (h *Handler) Calculate(c *gin.Context) {
	var req calculationRequest
	if !httpx.BindJSON(c, &req) {
		// BindJSON already wrote the error response.
		return
	}

	// Pre-validate length and charset (handler layer; returns validation_failed).
	if len(req.Expression) > 256 {
		h.logger.InfoContext(c.Request.Context(), "calculate validation failed",
			slog.String("reason", "expression too long"),
			slog.Int("length", len(req.Expression)),
		)
		httpx.WriteError(c, http.StatusBadRequest, httpx.CodeValidation,
			"invalid expression",
			map[string]string{"expression": "expression exceeds maximum length of 256 characters"})
		return
	}

	if !allowedCharsRE.MatchString(req.Expression) {
		h.logger.InfoContext(c.Request.Context(), "calculate validation failed",
			slog.String("reason", "illegal characters"),
			slog.String("expression", req.Expression),
		)
		httpx.WriteError(c, http.StatusBadRequest, httpx.CodeValidation,
			"invalid expression",
			map[string]string{"expression": "expression contains illegal characters"})
		return
	}

	result, err := h.evaluator(req.Expression)
	if err != nil {
		h.handleCalcError(c, req.Expression, err)
		return
	}

	c.JSON(http.StatusOK, calculationResult{
		Result:     result,
		Expression: req.Expression,
	})
}

// handleCalcError maps calc sentinel errors to the wire error codes.
func (h *Handler) handleCalcError(c *gin.Context, expr string, err error) {
	switch {
	case errors.Is(err, calc.ErrDivideByZero):
		h.logger.InfoContext(c.Request.Context(), "calculate error",
			slog.String("reason", "division by zero"),
			slog.String("expression", expr),
		)
		httpx.WriteError(c, http.StatusUnprocessableEntity, httpx.CodeCalculation,
			"division by zero", nil)

	case errors.Is(err, calc.ErrNonFinite):
		h.logger.InfoContext(c.Request.Context(), "calculate error",
			slog.String("reason", "non-finite result"),
			slog.String("expression", expr),
		)
		httpx.WriteError(c, http.StatusUnprocessableEntity, httpx.CodeCalculation,
			"result is not a finite number", nil)

	case errors.Is(err, calc.ErrInvalidExpression):
		h.logger.InfoContext(c.Request.Context(), "calculate validation failed",
			slog.String("reason", "invalid expression"),
			slog.String("expression", expr),
			slog.String("error", err.Error()),
		)
		httpx.WriteError(c, http.StatusBadRequest, httpx.CodeValidation,
			"invalid expression",
			map[string]string{"expression": err.Error()})

	default:
		h.logger.ErrorContext(c.Request.Context(), "calculate unexpected error",
			slog.String("expression", expr),
			slog.String("error", err.Error()),
		)
		httpx.WriteError(c, http.StatusInternalServerError, httpx.CodeInternal,
			"internal server error", nil)
	}
}
