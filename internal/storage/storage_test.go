package storage_test

import (
	"encoding/json"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	t.Parallel()

	// Создаем временный файл для тестирования.
	tempFile, err := os.CreateTemp("", "testfilestorage")
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удалить временный файл после завершения теста.

	// Создаем экземпляр InMemoryStorage и передаем ему созданный временный файл.
	db := storage.NewInMemoryStorage()

	// Добавляем данные в хранилище.
	alias := "2d264"
	url := "https://mail.ru"
	err = db.Save(tempFile.Name(), alias, url)

	assert.Nil(t, err)

	// Открываем файл для чтения.
	file, err := os.Open(tempFile.Name())
	assert.Nil(t, err)
	defer file.Close()

	// Читаем данные из файла.
	var dataInFile struct {
		Alias string `json:"short_url"`
		URL   string `json:"original_url"`
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dataInFile)
	assert.Nil(t, err)

	// Проверяем, что данные в файле соответствуют ожидаемым данным.
	assert.Equal(t, "2d264", dataInFile.Alias)
	assert.Equal(t, "https://mail.ru", dataInFile.URL)
}
