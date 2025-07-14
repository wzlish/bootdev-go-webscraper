package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		expectError bool
		error       error
	}{
		{
			name:        "remove scheme",
			inputURL:    "https://blog.boot.dev/path",
			expected:    "blog.boot.dev/path",
			expectError: false,
		},
		{
			name:        "longer path",
			inputURL:    "https://blog.boot.dev/a/b/c",
			expected:    "blog.boot.dev/a/b/c",
			expectError: false,
		},
		{
			name:        "Invalid scheme",
			inputURL:    "mailto:notanaddress@null.void",
			expected:    "",
			expectError: true,
			error:       ErrInvalidScheme,
		},
		{
			name:        "Query string",
			inputURL:    "https://blog.boot.dev/?salmon=fiend",
			expected:    "blog.boot.dev/?salmon=fiend",
			expectError: false,
		},
		{
			name:        "Query string with path",
			inputURL:    "https://blog.boot.dev/a/b/c/d/?salmon=fiend",
			expected:    "blog.boot.dev/a/b/c/d/?salmon=fiend",
			expectError: false,
		},
		{
			name:        "Empty string",
			inputURL:    "https://",
			expected:    "",
			expectError: true,
			error:       ErrInvalidHost,
		},
		{
			name:        "Ip address host",
			inputURL:    "http://127.0.0.1/a/b/c/d",
			expected:    "127.0.0.1/a/b/c/d",
			expectError: false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				if !tc.expectError {
					t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
					return
				}
				if !errors.Is(err, tc.error) {

					fmt.Printf("Got: %v\n", err)
					fmt.Printf("Expected: %v\n", tc.error)

					t.Errorf("Test %v - '%s' FAIL: expected %v got %v", i, tc.name, tc.error, err)
					return
				}
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
