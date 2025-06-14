package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fieldName := strings.ToLower(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				errors[fieldName] = fmt.Sprintf("%s is required", fieldErr.Field())
			case "email":
				errors[fieldName] = fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
			case "min":
				errors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
			case "max":
				errors[fieldName] = fmt.Sprintf("%s must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
			default:
				errors[fieldName] = fmt.Sprintf("%s is not valid", fieldErr.Field())
			}
		}
	}
	return errors
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
