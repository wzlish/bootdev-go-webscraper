package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestGetHTML(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		inputBody   string
		expected    []string
		expectError bool
		error       error // Expected error type for specific error checks
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
</html>
`,
			expected:    []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
			expectError: false,
		},
		{
			name:        "empty HTML body",
			inputURL:    "https://example.com",
			inputBody:   "",
			expected:    []string{},
			expectError: false, // html.Parse handles empty string without error
		},
		{
			name:        "HTML body with no a tags",
			inputURL:    "https://example.com",
			inputBody:   `<html><body><p>No links here</p></body></html>`,
			expected:    []string{},
			expectError: false,
		},
		{
			name:        "HTML body with a tags but no href attribute",
			inputURL:    "https://example.com",
			inputBody:   `<html><body><a>Link text</a></body></html>`,
			expected:    []string{},
			expectError: false,
		},
		{
			name:     "HTML body with mixed valid and invalid hrefs",
			inputURL: "https://example.com",
			inputBody: `
<html>
			<body>
				<a href="/valid-path">Valid Link</a>
				<a href="invalid- url">Invalid Link (space)</a>
				<a href="http://another.com/path">Another Valid Link</a>
				<a href="javascript:alert('xss')">JS Link</a>
			</body>
</html>
`,
			// The function logs errors for malformed hrefs but still returns valid ones: note the automatic space to %20
			expected:    []string{"https://example.com/valid-path", "https://example.com/invalid-%20url", "http://another.com/path"},
			expectError: false,
		},
		{
			name:     "HTML body with duplicate URLs",
			inputURL: "https://example.com",
			inputBody: `
<html>
			<body>
				<a href="/path/duplicate">Link 1</a>
				<a href="https://example.com/path/duplicate">Link 2</a>
				<a href="/another-path">Link 3</a>
			</body>
</html>
`,
			// The function does not de-duplicate, so both will appear.
			expected:    []string{"https://example.com/path/duplicate", "https://example.com/path/duplicate", "https://example.com/another-path"},
			expectError: false,
		},
		{
			name:     "HTML body with relative URLs resolving to same absolute URL",
			inputURL: "https://example.com/dir/",
			inputBody: `
<html>
			<body>
				<a href="../path/one">Link 1</a>
				<a href="/path/one">Link 2</a>
			</body>
</html>
`,
			expected:    []string{"https://example.com/path/one", "https://example.com/path/one"},
			expectError: false,
		},
		{
			name:        "invalid rawBaseURL",
			inputURL:    "invalid url",
			inputBody:   `<html><body><a href="/path">Link</a></body></html>`,
			expected:    nil,
			expectError: true,
			error:       ErrInvalidBaseURL,
		},
		{
			name:     "HTML with fragment and query parameters",
			inputURL: "https://example.com/page",
			inputBody: `
<html>
			<body>
				<a href="/section#about">About Section</a>
				<a href="?param=value&another=true">Query Link</a>
				<a href="full/path.html?id=123#top">Full Link</a>
			</body>
</html>
`,
			expected:    []string{"https://example.com/section#about", "https://example.com/page?param=value&another=true", "https://example.com/full/path.html?id=123#top"},
			expectError: false,
		},
		{
			name:        "HTML with self-referencing href",
			inputURL:    "https://example.com/current/page.html",
			inputBody:   `<html><body><a href="">Self Link</a><a href=".">Dot Link</a><a href="./">Dot Slash Link</a></body></html>`,
			expected:    []string{"https://example.com/current/page.html", "https://example.com/current/", "https://example.com/current/"}, // Note: "./" resolves to the directory
			expectError: false,
		},
		{
			name:     "HTML with external links only",
			inputURL: "https://example.com",
			inputBody: `
<html>
			<body>
				<a href="https://google.com">Google</a>
				<a href="http://bing.com">Bing</a>
			</body>
</html>
`,
			expected:    []string{"https://google.com", "http://bing.com"},
			expectError: false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
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
				return // error == error, so pass.
			}
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test %v - %s FAIL: expected HTML: %v, Actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
