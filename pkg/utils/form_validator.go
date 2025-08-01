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
				// Convert field name to match form field names
				fieldName := convertFieldName(e.Field())

				// Add the field name and its error message to the map
				switch e.Tag() {
				case "required":
					errors[fieldName] = "Preenchimento obrigatório"
				case "min":
					errors[fieldName] = fmt.Sprintf("Digite no mínimo %s caracteres", e.Param())
				case "max":
					errors[fieldName] = fmt.Sprintf("Digite no máximo %s caracteres", e.Param())
				case "gt":
					errors[fieldName] = "Valor deve ser maior que zero"
				case "email":
					errors[fieldName] = "Email inválido"
				default:
					errors[fieldName] = "Revise este campo"
				}

			}
			return errors
		}
	}
	return nil
}

// convertFieldName converts struct field names to form field names
func convertFieldName(fieldName string) string {
	switch fieldName {
	case "Title":
		return "title"
	case "Description":
		return "description"
	case "SalesPage":
		return "sales_page"
	case "Value":
		return "value"
	case "Status":
		return "status"
	default:
		return fieldName
	}
}
