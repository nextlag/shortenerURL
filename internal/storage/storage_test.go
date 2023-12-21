package storage_test

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/storage"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	log := zap.NewNop()
	defer os.Remove(fileStorage)
	data := storage.NewFileStorage("-", "12345", "https://yandex.ru")
	if err := storage.Save(log, fileStorage, data.Alias, data.URL); err != nil {
		t.Error(err)
	}
}
