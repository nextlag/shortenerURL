package usecase

import (
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	defer os.Remove(fileStorage)
	data := NewFileStorage("12345", "http://yandex.ru", 1)
	if err := Save(fileStorage, data.Alias, data.URL, 1); err != nil {
		t.Error(err)
	}
}
