package generatestring

import (
	"math/rand"

	"github.com/google/uuid"
)

// NewRandomString generates a random string of a given length.
func NewRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// GenerateUUID - generates a UUID.
func GenerateUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}
