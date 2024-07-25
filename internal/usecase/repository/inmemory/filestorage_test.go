package inmemory_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/usecase/repository/inmemory"
)

func TestNewFileStorage(t *testing.T) {
	fs := inmemory.NewFileStorage("243", "alias", "http://example.com")

	assert.Equal(t, "243", fs.UUID)
	assert.Equal(t, "alias", fs.Alias)
	assert.Equal(t, "http://example.com", fs.URL)
}

func TestProducer_WriteEvent(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName)

	producer, err := inmemory.NewProducer(fileName)
	assert.NoError(t, err)
	defer producer.Close()

	event := inmemory.NewFileStorage("243", "alias", "http://example.com")

	err = inmemory.WriteEvent(producer, event)
	assert.NoError(t, err)

	// Verify the file content
	file, err := os.Open(fileName)
	assert.NoError(t, err)
	defer file.Close()

	consumer, err := inmemory.NewConsumer(fileName)
	assert.NoError(t, err)
	defer consumer.Close()

	readEvent, err := inmemory.ReadEvent[inmemory.FileStorage](consumer)
	assert.NoError(t, err)

	assert.Equal(t, *event, readEvent)
}

func TestConsumer_ReadEvent(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName)

	event := inmemory.NewFileStorage("12345", "alias", "http://example.com")

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	assert.NoError(t, err)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(event)
	assert.NoError(t, err)
	file.Close()

	consumer, err := inmemory.NewConsumer(fileName)
	assert.NoError(t, err)
	defer consumer.Close()

	readEvent, err := inmemory.ReadEvent[inmemory.FileStorage](consumer)
	assert.NoError(t, err)
	assert.Equal(t, *event, readEvent)
}

func TestProducer_Close(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName)

	producer, err := inmemory.NewProducer(fileName)
	assert.NoError(t, err)

	err = producer.Close()
	assert.NoError(t, err)
}

func TestConsumer_Close(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName)

	consumer, err := inmemory.NewConsumer(fileName)
	assert.NoError(t, err)

	err = consumer.Close()
	assert.NoError(t, err)
}
