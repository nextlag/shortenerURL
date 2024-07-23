package generatestring

import (
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
)

// NewRandomString generates a random string of a given length using crypto/rand for better security.
func NewRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	charLen := big.NewInt(int64(len(chars)))
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, charLen)
		if err != nil {
			// Handle error
			return ""
		}
		result[i] = chars[index.Int64()]
	}
	return string(result)
}

// GenerateUUID - generates a UUID.
func GenerateUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}
