package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileRepository_CodeExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_codes.txt")

	content := "abc\nxyz\npromo\nsale\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	repo := NewFileCodeRepository(testFile)

	testCases := []struct {
		name     string
		code     string
		expected bool
	}{
		{"first line", "abc", true},
		{"middle", "promo", true},
		{"last line", "sale", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.Exists(tc.code)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("code=%q: expected %v, got %v", tc.code, tc.expected, result)
			}
		})
	}
}

func TestFileRepository_CodeNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_codes.txt")

	content := "abc\nxyz\npromo\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	repo := NewFileCodeRepository(testFile)

	testCases := []struct {
		name string
		code string
	}{
		{"nonexistent code", "invalid"},
		{"empty code", ""},
		{"partial match", "ab"},
		{"case sensitive", "ABC"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.Exists(tc.code)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != false {
				t.Errorf("code=%q: expected false, got %v", tc.code, result)
			}
		})
	}
}

func TestFileRepository_FileError(t *testing.T) {
	repo := NewFileCodeRepository("/nonexistent/file.txt")

	result, err := repo.Exists("promo")

	if err == nil {
		t.Error("expected error for nonexistent file")
	}

	if result != false {
		t.Errorf("expected false for nonexistent file, got %v", result)
	}
}

func TestFileRepository_LoadAll(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_codes.txt")

	content := "abc\nxyz\npromo\nsale\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	repo := NewFileCodeRepository(testFile)
	codeSet, err := repo.LoadAll()

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := []string{"abc", "xyz", "promo", "sale"}

	if len(codeSet) != len(expected) {
		t.Errorf("expected %d codes, got %d", len(expected), len(codeSet))
	}

	for _, code := range expected {
		if _, exists := codeSet[code]; !exists {
			t.Errorf("expected code %q to be in set", code)
		}
	}
}
