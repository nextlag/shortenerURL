package generatestring

import (
	"github.com/google/uuid"
	"math/rand"
)

// NewRandomString генерирует случайную строку заданной длины.
func NewRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func GenerateUUID() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}
