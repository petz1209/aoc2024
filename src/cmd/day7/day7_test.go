package main

import "testing"

func TestDay7P1(t *testing.T) {
	expected := 3749
	got := puzzle1("test_input.txt")
	if expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}

}

func TestDay7P2(t *testing.T) {
	expected := 11387
	got := puzzle2("test_input.txt")
	if expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}

}
