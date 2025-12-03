package main

import (
	"errors"
	"fmt"
)

// Custom error type
type ValidationError struct {
	Field string
	Value interface{}
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' with value '%v' - %s", e.Field, e.Value, e.Msg)
}

// Sentinel errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
)

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{
			Field: "age",
			Value: age,
			Msg:   "age cannot be negative",
		}
	}
	if age > 150 {
		return &ValidationError{
			Field: "age",
			Value: age,
			Msg:   "age is unrealistic",
		}
	}
	return nil
}

func main() {
	// Error handling examples
	err := validateAge(-5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		// Type assertion
		var validErr *ValidationError
		if errors.As(err, &validErr) {
			fmt.Printf("Field: %s\n", validErr.Field)
		}
	}

	// Error wrapping (Go 1.13+)
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("wrapped: %w", baseErr)

	fmt.Printf("\nWrapped error: %v\n", wrappedErr)
	fmt.Printf("Is base error: %t\n", errors.Is(wrappedErr, baseErr))
}
