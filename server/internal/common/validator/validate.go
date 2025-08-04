package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationErrorDetail struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value"`
}

func ValidateStruct(s any) ([]ValidationErrorDetail, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	var errs []ValidationErrorDetail
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errs = append(errs, ValidationErrorDetail{
				Field: strings.ToLower(e.Field()),
				Tag:   e.Tag(),
				Value: e.Value(),
			})
		}
		return errs, fmt.Errorf("validation failed")
	}

	return nil, err
}
