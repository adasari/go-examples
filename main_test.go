package main

import "testing"

func TestTry(t *testing.T) {
	got := try()

	if got != nil {
		t.Fatalf("failed to try: %v", got)
	}
}
