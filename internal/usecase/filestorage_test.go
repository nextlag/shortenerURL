package usecase_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/usecase"
)

func TestNewFileStorage(t *testing.T) {
	fs := usecase.NewFileStorage("12345", "http://example.com", "1")

	assert.Equal(t, "1", fs.UUID)
	assert.Equal(t, "12345", fs.Alias)
	assert.Equal(t, "http://example.com", fs.URL)
}

func TestProducer_WriteEvent(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName) // Clean up the file after the test

	producer, err := usecase.NewProducer(fileName)
	assert.NoError(t, err)
	defer producer.Close()

	event := &usecase.FileStorage{
		UUID:  "12345",
		Alias: "alias",
		URL:   "http://example.com",
	}

	err = producer.WriteEvent(event)
	assert.NoError(t, err)

	// Verify the file content
	file, err := os.Open(fileName)
	assert.NoError(t, err)
	defer file.Close()

	consumer, err := usecase.NewConsumer(fileName)
	assert.NoError(t, err)
	defer consumer.Close()

	readEvent, err := consumer.ReadEvent()
	assert.NoError(t, err)
	assert.Equal(t, event, readEvent)
}

func TestConsumer_ReadEvent(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName) // Clean up the file after the test

	// Write a test event to the file
	event := &usecase.FileStorage{
		UUID:  "12345",
		Alias: "alias",
		URL:   "http://example.com",
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	assert.NoError(t, err)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(event)
	assert.NoError(t, err)
	file.Close()

	// Test reading the event
	consumer, err := usecase.NewConsumer(fileName)
	assert.NoError(t, err)
	defer consumer.Close()

	readEvent, err := consumer.ReadEvent()
	assert.NoError(t, err)
	assert.Equal(t, event, readEvent)
}

func TestProducer_Close(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName) // Clean up the file after the test

	producer, err := usecase.NewProducer(fileName)
	assert.NoError(t, err)

	err = producer.Close()
	assert.NoError(t, err)
}

func TestConsumer_Close(t *testing.T) {
	fileName := "test_file_storage.json"
	defer os.Remove(fileName) // Clean up the file after the test

	consumer, err := usecase.NewConsumer(fileName)
	assert.NoError(t, err)

	err = consumer.Close()
	assert.NoError(t, err)
}
