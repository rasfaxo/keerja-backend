package utils

import (
	stdErrors "errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the validator
func InitValidator() {
	validate = validator.New()

	// Register custom validations
	registerCustomValidations()
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

func ValidateStruct(s interface{}) error {
	v := GetValidator()
	return v.Struct(s)
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if !stdErrors.As(err, &validationErrors) {
		return errors
	}

	for _, err := range validationErrors {
		field := strings.ToLower(err.Field())
		errors[field] = getErrorMessage(err)
	}

	return errors
}

func getErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func registerCustomValidations() {
	// Add custom validations here if needed
	// Example: validate.RegisterValidation("custom_tag", customValidationFunc)
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPhone(phone string) bool {
	// Indonesian phone number validation
	// Accepts: 08xxxxxxxxxx, +628xxxxxxxxxx, 628xxxxxxxxxx
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)[0-9]{9,12}$`)
	return phoneRegex.MatchString(phone)
}

func IsValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(:[0-9]+)?(/.*)?$`)
	return urlRegex.MatchString(url)
}

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return stdErrors.New("password must be at least 8 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper {
		return stdErrors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return stdErrors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return stdErrors.New("password must contain at least one number")
	}

	return nil
}

func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return email
}
