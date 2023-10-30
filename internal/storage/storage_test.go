package storage_test

import (
	"os"
	"testing"

	"github.com/nextlag/shortenerURL/internal/lib/filestorage"
	"github.com/nextlag/shortenerURL/internal/storage"
)

func TestSettings(t *testing.T) {
	fname := "file_test.json"
	defer os.Remove(fname)
	data := filestorage.Request{
		URL:   "http://yandex.ru",
		Alias: "12345",
	}
	if err := storage.Save(fname, data.Alias, data.URL); err != nil {
		t.Error(err)
	}
}
