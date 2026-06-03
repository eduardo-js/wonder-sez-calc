package httpx

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Register a TagNameFunc so validator uses the json struct tag as the field
	// name in ValidationErrors, giving callers json-friendly field keys.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
	}
}

// BindJSON decodes and validates the request body into dst.
//
// Returns true on success. On failure it writes the standard error envelope,
// aborts the context, and returns false.
//
//   - validator.ValidationErrors  -> 400 CodeValidation with per-field details
//   - any other decode/syntax error -> 400 CodeBadRequest
func BindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string, len(ve))
			for _, fe := range ve {
				fields[fe.Field()] = fe.Tag()
			}
			WriteError(c, http.StatusBadRequest, CodeValidation, "validation failed", fields)
			return false
		}
		WriteError(c, http.StatusBadRequest, CodeBadRequest, "bad request", nil)
		return false
	}
	return true
}
