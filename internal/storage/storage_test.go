package storage_test

import (
	"os"
	"testing"
	"time"

	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

func TestSettings(t *testing.T) {
	fname := "file_test.json"
	defer os.Remove(fname)
	data := usecase.NewRequest(0, "12345", "https://yandex.ru", "-", time.Now())
	if err := storage.Save(fname, data.GetEntityRequest().Alias, data.GetEntityRequest().URL); err != nil {
		t.Error(err)
	}
}
