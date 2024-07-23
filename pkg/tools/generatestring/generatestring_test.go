package generatestring

import (
	"testing"
	"unicode/utf8"

	"github.com/google/uuid"
)

func TestNewRandomString(t *testing.T) {
	length := 10
	randomString := NewRandomString(length)

	// Check if the generated string has the correct length
	if len(randomString) != length {
		t.Errorf("expected length %d, but got %d", length, len(randomString))
	}

	// Check if the string contains only allowed characters
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range randomString {
		if !utf8.ValidRune(char) || !containsRune(allowedChars, char) {
			t.Errorf("generated string contains invalid character: %c", char)
		}
	}
}

func containsRune(s string, r rune) bool {
	for _, v := range s {
		if v == r {
			return true
		}
	}
	return false
}

func TestGenerateUUID(t *testing.T) {
	uuidStr := GenerateUUID()

	// Check if the generated UUID has the correct length
	expectedLength := 36
	if len(uuidStr) != expectedLength {
		t.Errorf("expected UUID length %d, but got %d", expectedLength, len(uuidStr))
	}

	// Check if the generated UUID is valid
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("generated UUID is invalid: %s", uuidStr)
	}
}
