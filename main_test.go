package main

import "testing"

func TestFormatSize(t *testing.T) {
	// define table of test cases
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Bytes", 500, "500 B"},
		{"Exact Kilobyte", 1024, "1.0 KB"},
		{"Exact Megabyte", 1024 * 1024, "1.0 MB"},
		{"Some Kilobyte", 1536, "1.5 KB"},
		{"Some Megabyte", 1536 * 1024, "1.5 MB"},
	}

	// loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatSize(tt.input)

			if actual != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, actual)
			}
		})
	}
}
