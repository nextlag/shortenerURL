package usecase

import (
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	defer os.Remove(fileStorage)
	data := NewFileStorage("-", "12345", "https://yandex.ru")
	if err := Save(fileStorage, data.Alias, data.URL); err != nil {
		t.Error(err)
	}
}
