package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateHistoryTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectDate  string
	}{
		{"valid january", "202501", false, "2025/01/01"},
		{"valid december", "199912", false, "1999/12/01"},
		{"invalid month", "202513", true, ""},
		{"wrong length", "20259", true, ""},
		{"non-numeric", "abcd01", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateHistoryTime(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for input %q, got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error for input %q: %v", tt.input, err)
				return
			}

			want, _ := time.Parse("2006/01/02", tt.expectDate)
			if !got.Equal(want) {
				t.Errorf("input %q -> got %v, want %v", tt.input, got, want)
			}
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	tmpDir := t.TempDir()

	validPath := filepath.Join(tmpDir, "out.csv")
	if err := os.WriteFile(validPath, []byte("date,amount,content\n"), 0644); err != nil {
		t.Fatalf("failed to create valid file: %v", err)
	}

	dirPath := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectedExt string
		expectError bool
	}{
		{"valid csv", validPath, ".csv", false},
		{"wrong extension", validPath, ".json", true},
		{"directory instead of file", dirPath, ".csv", true},
		{"nonexistent file", filepath.Join(tmpDir, "missing.csv"), ".csv", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFilePath(tt.path, tt.expectedExt)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil (path=%s, ext=%s)", tt.path, tt.expectedExt)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v (path=%s, ext=%s)", err, tt.path, tt.expectedExt)
				}
			}
		})
	}
}
