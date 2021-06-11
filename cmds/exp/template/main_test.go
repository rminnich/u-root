package main

import "testing"

func TestEmpty(t *testing.T) {
	s, err := template()
	if err != nil {
		t.Fatalf("No template: got %v, want nil", err)
	}
	if len(s) != 0 {
		t.Fatalf("No template: got %d strings, want 0", len(s))
	}
}

func TestOverlap(t *testing.T) {
	s, err := template("core", "core", "all")
	if err != nil {
		t.Fatalf("No template: got %v, want nil", err)
	}
	if len(s) != 0 {
		t.Fatalf("No template: got %d strings, want 0", len(s))
	}
}
