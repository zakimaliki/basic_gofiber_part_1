package helpers

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(param any) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(param)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := ErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}
