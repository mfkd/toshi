package errors

import (
	"errors"
	"fmt"
)

// Wrap wraps an error with additional context.
func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// Unwrap unwraps an error to retrieve its cause.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is checks if an error is of a specific type.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// New creates a new error.
func New(msg string) error {
	return errors.New(msg)
}
