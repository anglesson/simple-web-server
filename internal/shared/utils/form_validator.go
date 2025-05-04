package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateForm(form interface{}) map[string]string {
	validate := validator.New()
	err := validate.Struct(form)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Create a map to store field-specific error messages
			errors := make(map[string]string)
			for _, e := range validationErrors {
				// Add the field name and its error message to the map
				switch e.Tag() {
				case "required":
					errors[e.Field()] = "Preenchimento obrigatório"
				case "min":
					errors[e.Field()] = fmt.Sprintf("Digite no mínimo %s caracteres", e.Param())
				case "max":
					errors[e.Field()] = fmt.Sprintf("Digite no máximo %s caracteres", e.Param())
				default:
					errors[e.Field()] = "Revise este campo"
				}

			}
			return errors
		}
	}
	return nil
}
