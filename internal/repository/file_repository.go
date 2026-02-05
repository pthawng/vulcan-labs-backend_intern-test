package repository

import (
	"bufio"
	"os"
)

// FileCodeRepository implements CodeRepository using a file.
type FileCodeRepository struct {
	filePath string
}

func NewFileCodeRepository(filePath string) *FileCodeRepository {
	return &FileCodeRepository{filePath: filePath}
}

func (r *FileCodeRepository) Exists(code string) (bool, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == code {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// LoadAll loads all codes from the file into a set for O(1) lookup.
// Uses map[string]struct{} for zero memory overhead on values.
// This is typically called once and cached by the service layer.
func (r *FileCodeRepository) LoadAll() (map[string]struct{}, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	codeSet := make(map[string]struct{})
	scanner := bufio.NewScanner(file)

	// Tune buffer size for future-proofing
	// Default: 64KB max token size
	// Increased to 1MB to handle potential longer codes in the future
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		code := scanner.Text()
		codeSet[code] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return codeSet, nil
}
