package testutil

import (
	"reflect"
	"testing"
)

// Equal is a helper function to check equality of two values.
func Equal(t *testing.T, expected, actual any, optionalMessage ...string) {
	t.Helper()

	var message string
	if len(optionalMessage) > 0 {
		message = optionalMessage[0]
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v %q", expected, actual, message)
	}
}

// True is Helper function to check if a value is true.
func True(t *testing.T, actual bool, optionalMessage ...string) {
	t.Helper()

	var message string
	if len(optionalMessage) > 0 {
		message = optionalMessage[0]
	}

	if !actual {
		t.Fatal("expected true, got false", message)
	}
}

// EqualSlices is a helper function to check equality of two slices.
func EqualSlices[T1, T2 any](t *testing.T, expected []T1, actual []T2) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("expected slice length %d, got %d", len(expected), len(actual))

		return
	}

	for i := range expected {
		if !reflect.DeepEqual(expected[i], actual[i]) {
			t.Fatalf("at index %d: expected %v, got %v", i, expected[i], actual[i])
		}
	}
}

// Panics is a helper function to check if a function panics.
func Panics(t *testing.T, fnc func(), optionalMessage ...string) {
	t.Helper()

	var message string
	if len(optionalMessage) > 0 {
		message = optionalMessage[0]
	}

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic, got nil %q", message)
		}
	}()

	fnc()
}
