package main

import "testing"

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	if result != expected {
		t.Fatalf("Add(2, 3) = %d; expected %d", result, expected)
	}
}

func TestSubtract(t *testing.T) {
	result := Subtract(2, 2)
	expected := 0
	if result != expected {
		t.Fatalf("Subtract(2, 2) = %d; expected %d", result, expected)
	}
}
