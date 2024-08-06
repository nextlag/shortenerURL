package inmemory

import (
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	defer os.Remove(fileStorage)
	defer os.Remove(fileDel)
	data := NewFileStorage("1", "12345", "http://yandex.ru")
	if err := save(fileStorage, data.Alias, data.URL, 1, false); err != nil {
		t.Error(err)
	}
}