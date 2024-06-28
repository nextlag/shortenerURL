package generatestring_test

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

//go:generate godoc -http=:8090 -play

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

// Example demonstrates how to use functions from the generatestring package.
func Example() {
	randomString := NewRandomString(10)
	fmt.Println("Random String:", randomString)

	uuid := GenerateUUID()
	fmt.Println("UUID:", uuid)
}
