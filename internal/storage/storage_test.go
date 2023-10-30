package storage_test

import (
	"github.com/nextlag/shortenerURL/internal/lib/storagefile"
	"github.com/nextlag/shortenerURL/internal/storage"
	"os"
	"testing"
)

func TestSettings(t *testing.T) {
	fname := "file_test.json"
	defer os.Remove(fname)
	data := storagefile.Event{
		URL:   "http://yandex.ru",
		Alias: "12345",
	}
	if err := storage.Save(fname, data.Alias, data.URL); err != nil {
		t.Error(err)
	}
}
