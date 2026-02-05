package promotion

import (
	"testing"
)

func TestValidateCode_Valid(t *testing.T) {
	err := ValidateCode("promo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCode_Empty(t *testing.T) {
	err := ValidateCode("")

	if err == nil {
		t.Error("expected error for empty code")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	expectedMsg := "code cannot be empty"
	if validationErr.Message != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, validationErr.Message)
	}
}

func TestValidateCode_TooLong(t *testing.T) {
	err := ValidateCode("abcdef")

	if err == nil {
		t.Error("expected error for code too long")
		return
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	expectedMsg := "code must be at most 5 characters"
	if validationErr.Message != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, validationErr.Message)
	}
}

func TestValidateCode_InvalidCharacters(t *testing.T) {
	testCases := []struct {
		name string
		code string
	}{
		{"uppercase", "Promo"},
		{"numbers", "abc12"},
		{"special chars", "ab_cd"},
		{"space", "ab cd"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCode(tc.code)

			if err == nil {
				t.Errorf("expected error for %s", tc.name)
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			expectedMsg := "code must contain only lowercase letters (a-z)"
			if validationErr.Message != expectedMsg {
				t.Errorf("expected error message %q, got %q", expectedMsg, validationErr.Message)
			}
		})
	}
}
