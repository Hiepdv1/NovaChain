package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationErrorDetail struct {
	Field        string `json:"field"`
	Tag          string `json:"tag"`
	Param        string `json:"param,omitempty"`
	Expected     string `json:"expected"`
	CurrentType  string `json:"current_type"`
	CurrentValue any    `json:"current_value"`
	Message      string `json:"message"`
}

func ValidateStruct(s any) ([]ValidationErrorDetail, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	var errs []ValidationErrorDetail
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldPath := strings.ToLower(e.StructNamespace())
			fieldPath = trimRoot(fieldPath)

			kind := e.Kind()
			currentType := e.Type().String()
			param := e.Param()

			expected := buildExpectedMessage(e.Tag(), param, kind)
			message := fmt.Sprintf("Field '%s' %s. Got: %v (%s)", fieldPath, expected, e.Value(), currentType)

			errs = append(errs, ValidationErrorDetail{
				Field:        fieldPath,
				Tag:          e.Tag(),
				Param:        param,
				Expected:     expected,
				CurrentType:  currentType,
				CurrentValue: e.Value(),
				Message:      message,
			})
		}
		return errs, fmt.Errorf("validation failed")
	}

	return nil, err
}

func trimRoot(ns string) string {
	if idx := strings.Index(ns, "."); idx >= 0 {
		return ns[idx+1:]
	}
	return ns
}

func buildExpectedMessage(tag, param string, kind reflect.Kind) string {
	switch tag {
	case "required":
		return "is required (must not be empty)"
	case "len":
		if isStringKind(kind) {
			return fmt.Sprintf("must have length = %s characters", param)
		}
		return fmt.Sprintf("must have length = %s items", param)
	case "min":
		if isNumericKind(kind) {
			return fmt.Sprintf("must be >= %s", param)
		}
		return fmt.Sprintf("length must be >= %s", param)
	case "max":
		if isNumericKind(kind) {
			return fmt.Sprintf("must be <= %s", param)
		}
		return fmt.Sprintf("length must be <= %s", param)
	case "gt":
		return fmt.Sprintf("must be > %s", param)
	case "gte":
		return fmt.Sprintf("must be >= %s", param)
	case "lt":
		return fmt.Sprintf("must be < %s", param)
	case "lte":
		return fmt.Sprintf("must be <= %s", param)
	case "eq":
		return fmt.Sprintf("must be equal to %s", param)
	case "ne":
		return fmt.Sprintf("must not be equal to %s", param)
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", strings.ReplaceAll(param, " ", ", "))
	case "email":
		return "must be a valid email address"
	case "uuid", "uuid4":
		return "must be a valid UUID"
	case "hexadecimal", "hex":
		return "must be a hexadecimal string"
	case "alphanum":
		return "must contain only alphanumeric characters"
	case "alpha":
		return "must contain only letters"
	case "numeric":
		return "must be a numeric value"
	case "boolean":
		return "must be a boolean value"
	case "required_if":
		return fmt.Sprintf("is required when %s", param)
	case "omitempty":
		return "optional (empty is allowed)"
	default:
		if param != "" {
			return fmt.Sprintf("%s=%s", tag, param)
		}
		return tag
	}
}

func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func isStringKind(k reflect.Kind) bool {
	return k == reflect.String
}
