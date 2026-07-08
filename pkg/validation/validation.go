package validation

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.Split(fld.Tag.Get("json"), ",")[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
	}
}

func FormatValidationErrors(ve validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, e := range ve {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errs[field] = field + " is required"
		case "min":
			errs[field] = field + " must be greater than " + e.Param()
		case "max":
			errs[field] = field + " must be less than " + e.Param()
		default:
			errs[field] = field + " is invalid"
		}
	}
	return errs
}
