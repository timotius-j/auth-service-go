package util

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct based on `validate` tags.
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// ValidatePhone validates phone number
func ValidatePhone(num string) bool {
	if len(num) < 10 {
		return false
	}
	if num[0] != '+' {
		return false
	}

	digits := num[1:]
	if len(digits) > 15 {
		return false
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return false
		}
	}
	return true
}
