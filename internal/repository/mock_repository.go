package repository

// MockCodeRepository is a mock implementation for testing.
type MockCodeRepository struct {
	codes map[string]bool
	err   error
}

// NewMockCodeRepository creates a new mock repository.
func NewMockCodeRepository(codes []string, err error) *MockCodeRepository {
	codeMap := make(map[string]bool)
	for _, code := range codes {
		codeMap[code] = true
	}
	return &MockCodeRepository{
		codes: codeMap,
		err:   err,
	}
}

// Exists checks if a code exists in the mock repository.
func (m *MockCodeRepository) Exists(code string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.codes[code], nil
}

// LoadAll returns all codes as a set.
// Returns a copy to prevent external modification.
func (m *MockCodeRepository) LoadAll() (map[string]struct{}, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Convert internal map[string]bool to map[string]struct{}
	result := make(map[string]struct{}, len(m.codes))
	for code := range m.codes {
		result[code] = struct{}{}
	}

	return result, nil
}
