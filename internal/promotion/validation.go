package promotion

import "fmt"

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// ValidateCode validates a promotion code according to business rules.
func ValidateCode(code string) error {
	if len(code) == 0 {
		return &ValidationError{
			Field:   "code",
			Message: "code cannot be empty",
		}
	}

	if len(code) > 5 {
		return &ValidationError{
			Field:   "code",
			Message: "code must be at most 5 characters",
		}
	}

	for _, char := range code {
		if char < 'a' || char > 'z' {
			return &ValidationError{
				Field:   "code",
				Message: "code must contain only lowercase letters (a-z)",
			}
		}
	}

	return nil
}
