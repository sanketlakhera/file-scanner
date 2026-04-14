package main

import "testing"

func TestFormatSize(t *testing.T) {
	expected := "500 B"
	actual := formatSize(500)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
