package main

import "testing"

func TestDay5Puzzle1(t *testing.T) {
	expect := 41
	got := solvePuzzle1("test_input.txt")

	if expect != got {
		t.Errorf("expected %d  but got %d", expect, got)
	}
}

func TestDay5Puzzle2(t *testing.T) {
	expect := 6
	got := solvePuzzle2("test_input.txt")

	if expect != got {
		t.Errorf("expected %d  but got %d", expect, got)
	}
}
