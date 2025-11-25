package utils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/herman-xphp/my-notes-api/internal/service"
)

// WrapError wraps an error with additional context
func WrapError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// WrapErrorf wraps an error with formatted context
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// HandleDBError converts database errors to service errors
func HandleDBError(err error, notFoundErr error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFoundErr
	}

	return err
}

// IfError executes function if error is not nil
func IfError(err error, fn func(error) error) error {
	if err != nil {
		return fn(err)
	}
	return nil
}

// ReturnIfError returns immediately if error is not nil
// Useful for reducing nested if statements
func ReturnIfError(err error, wrapper func(error) error) error {
	if err == nil {
		return nil
	}
	if wrapper != nil {
		return wrapper(err)
	}
	return err
}

// Must panics if error is not nil (use only in initialization)
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustValue returns value or panics if error (use only in initialization)
func MustValue[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, service.ErrUserNotFound) ||
		errors.Is(err, service.ErrNoteNotFound) ||
		errors.Is(err, gorm.ErrRecordNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, service.ErrEmailAlreadyExists) ||
		errors.Is(err, service.ErrInvalidCredentials)
}

// CheckAndWrap checks error and wraps it with context
func CheckAndWrap(err error, operation string) error {
	if err == nil {
		return nil
	}

	// Don't wrap business logic errors
	if IsValidationError(err) || IsNotFoundError(err) {
		return err
	}

	return WrapError(operation, err)
}
