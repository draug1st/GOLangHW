package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ReadFileTrim reads the named file and returns the trimmed content.
// It wraps any error encountered during reading.
func ReadFileTrim(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func main() {
	// Create a temporary file for success case
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	content := "   hello world   \n"
	if _, err := tmpFile.WriteString(content); err != nil {
		panic(err)
	}
	tmpFile.Close()

	// Case 1: Success
	got, err := ReadFileTrim(tmpFile.Name())
	if err != nil {
		fmt.Printf("Case 1 FAIL: unexpected error: %v\n", err)
	} else if got != "hello world" {
		fmt.Printf("Case 1 FAIL: expected 'hello world', got %q\n", got)
	} else {
		fmt.Println("Case 1 PASS: file read and trimmed correctly")
	}

	// Case 2: Error wrapping
	nonExistentPath := "non_existent_file.txt"
	_, err = ReadFileTrim(nonExistentPath)
	if err == nil {
		fmt.Println("Case 2 FAIL: expected error, got nil")
		return
	}

	// Check if the error string contains the context
	expectedMsgPart := fmt.Sprintf("read %q", nonExistentPath)
	if !strings.Contains(err.Error(), expectedMsgPart) {
		fmt.Printf("Case 2 FAIL: error message %q does not contain %q\n", err.Error(), expectedMsgPart)
	} else {
		fmt.Println("Case 2 PASS: error message contains context")
	}

	// Check if errors.As works (Unwrap)
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		fmt.Println("Case 3 PASS: errors.As found *os.PathError")
	} else {
		fmt.Printf("Case 3 FAIL: errors.As did not find *os.PathError in %v\n", err)
	}
}
