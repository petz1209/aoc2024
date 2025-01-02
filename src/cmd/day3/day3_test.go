package main

import (
	"reflect"
	"testing"
)

const INP = "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"

func TestContainer(t *testing.T) {
	got := extractNoise(INP)
	want := []string{"mul(2,4)", "mul(5,5)", "mul(11,8)", "mul(8,5)"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestExecutor(t *testing.T) {

	got := executor("mul(2,4)")
	want := 8
	if got != want {
		t.Errorf("got %d but wanted %d", got, want)
	}
}
