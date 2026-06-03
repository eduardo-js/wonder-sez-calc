package calculate

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Register mounts the calculate endpoint on the given router group using the
// default calc.Evaluate implementation.
//
//	POST /calculate
func Register(rg *gin.RouterGroup, logger *slog.Logger) {
	h := newHandler(logger)
	rg.POST("/calculate", h.Calculate)
}

// RegisterWithEvaluator mounts the calculate endpoint with a custom evaluator.
// Intended for testing only.
func RegisterWithEvaluator(rg *gin.RouterGroup, logger *slog.Logger, eval evaluatorFunc) {
	h := &Handler{logger: logger, evaluator: eval}
	rg.POST("/calculate", h.Calculate)
}
