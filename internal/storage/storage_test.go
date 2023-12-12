package storage_test

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

func TestSettings(t *testing.T) {
	fileStorage := "file_test.json"
	log := zap.NewNop()
	defer os.Remove(fileStorage)
	data := usecase.NewRequest("-", "12345", "https://yandex.ru")
	if err := storage.Save(log, fileStorage, data.GetEntityRequest().Alias, data.GetEntityRequest().URL); err != nil {
		t.Error(err)
	}
}
