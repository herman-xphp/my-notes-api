package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a custom validator wrapper
type Validator struct {
	validator *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()
	return &Validator{validator: v}
}

// Validate validates a struct
func (v *Validator) Validate(data any) error {
	if err := v.validator.Struct(data); err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors formats validation errors into readable messages
func (v *Validator) formatValidationErrors(err error) error {
	var errors []string

	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	for _, err := range validationError {
		errors = append(errors, v.formatFieldError(err))
	}

	return fmt.Errorf("%s", strings.Join(errors, "; "))
}

// formatFieldError formats a single field erros
func (v *Validator) formatFieldError(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be not exceed %s characters", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// GetValidationErrors extracts validation errors as a slice
func (v *Validator) GetValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	if err == nil {
		return validationErrors
	}

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		validationErrors = append(validationErrors, ValidationError{
			Field:   "",
			Message: err.Error(),
		})
		return validationErrors
	}

	for _, err := range valErrors {
		validationErrors = append(validationErrors, ValidationError{
			Field:   strings.ToLower(err.Field()),
			Message: v.formatFieldError(err),
		})
	}

	return validationErrors
}
