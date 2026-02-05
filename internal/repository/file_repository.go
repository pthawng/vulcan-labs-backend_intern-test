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
