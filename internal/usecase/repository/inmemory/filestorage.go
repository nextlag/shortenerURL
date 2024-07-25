package inmemory

import (
	"encoding/json"
	"os"
)

// FileStorage represents the structure of a URL record for file storage.
type FileStorage struct {
	UUID  string `json:"uuid"`                        // UUID, a unique identifier
	Alias string `json:"alias,omitempty"`             // Alias, a custom alias for the shortened URL (optional)
	URL   string `json:"url" validate:"required,url"` // URL, the URL to be shortened, must be a valid URL
}

// NewFileStorage creates a new instance of FileStorage.
func NewFileStorage(userID, alias, url string) *FileStorage {
	return &FileStorage{
		UUID:  userID,
		Alias: alias,
		URL:   url,
	}
}

// IsDeleted represents a deletion status of a URL record.
type IsDeleted struct {
	UserID    string `json:"uuid"`
	Alias     string `json:"alias"`
	StatusDel bool   `json:"status_del"`
}

// NewIsDeleted creates a new instance of IsDeleted.
func NewIsDeleted(userID, alias string, del bool) *IsDeleted {
	return &IsDeleted{
		UserID:    userID,
		Alias:     alias,
		StatusDel: del,
	}
}

// Producer is responsible for writing URL records to a file.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer creates a new Producer for the given file name.
func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteEvent writes a record to the file.
func WriteEvent[T any](p *Producer, event T) error {
	return p.encoder.Encode(event)
}

// Close closes the Producer's file.
func (p *Producer) Close() error {
	return p.file.Close()
}

// Consumer is responsible for reading records from a file.
type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer creates a new Consumer for the given file name.
func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadEvent reads a record from the file.
func ReadEvent[T any](c *Consumer) (T, error) {
	var event T
	if err := c.decoder.Decode(&event); err != nil {
		return event, err
	}
	return event, nil
}

// Close closes the Consumer's file.
func (c *Consumer) Close() error {
	return c.file.Close()
}
