package usecase

import (
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	fileDel := "file_del.json"
	defer os.Remove(fileDel)
	defer os.Remove(fileStorage)
	data := NewFileStorage("1", "12345", "http://yandex.ru")
	if err := save(fileStorage, fileDel, data.Alias, data.URL, 1, false); err != nil {
		t.Error(err)
	}
}
