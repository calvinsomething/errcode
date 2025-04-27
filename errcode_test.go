package main

import (
	"go/token"
	"testing"
)

func validateCode(t *testing.T, code string, expectedLength int) {
	if !token.IsIdentifier(code) || !token.IsExported(code) {
		t.Fatalf("code '%s' is not a valid exportable Go identifier", code)
	}

	if len(code) != expectedLength {
		t.Fatalf("invalid code length %d, expected %d", len(code), expectedLength)
	}
}

func TestNewCodeWithoutPrefix(t *testing.T) {
	prefix = ""

	for i := 0; i < 100; i++ {
		validateCode(t, newCode(), length)
	}
}

func TestNewCodeWithPrefix(t *testing.T) {
	prefix = "Test"

	for i := 0; i < 100; i++ {
		validateCode(t, newCode(), length+len(prefix))
	}
}
