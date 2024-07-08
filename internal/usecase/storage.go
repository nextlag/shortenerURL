// Package usecase provides use cases for managing short URLs,
// including in-memory data storage and file-based storage operations.
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/pkg/tools/generatestring"
)

// Data represents the in-memory data storage structure.
type Data struct {
	data  map[string]string
	log   *zap.Logger
	cfg   *configuration.Config
	mutex sync.Mutex // Mutex for synchronizing access to data
}

// NewData creates a new instance of Data.
func NewData(log *zap.Logger, cfg *configuration.Config) *Data {
	return &Data{
		data: make(map[string]string),
		log:  log,
		cfg:  cfg,
	}
}

// Get retrieves a URL by its alias from the in-memory storage.
func (s *Data) Get(_ context.Context, alias string) (string, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve from memory
	url, ok := s.data[alias]
	if !ok {
		return "", false, fmt.Errorf("key '%s' not found", alias)
	}
	return url, false, nil
}

// GetAll returns a message indicating that memory storage cannot operate with user IDs.
func (s *Data) GetAll(context.Context, int, string) ([]byte, error) {
	return []byte("Memory storage can't operate with user IDs"), nil
}

// Healthcheck always returns true for in-memory storage.
func (s *Data) Healthcheck() (bool, error) {
	return true, nil
}

// Del does nothing for in-memory storage as delete operations are not implemented.
func (s *Data) Del(_ int, _ []string) error {
	return nil
}

// Put saves a URL with a generated alias in the in-memory storage.
func (s *Data) Put(_ context.Context, url string, alias string, _ int) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if alias == "" {
		alias = generatestring.NewRandomString(8)
	}

	// Check for alias existence
	if _, exists := s.data[alias]; exists {
		return "", fmt.Errorf("alias '%s/%s' already exists", s.cfg.BaseURL, alias)
	}

	// Check for URL existence
	for k, v := range s.data {
		if v == url {
			return k, nil
		}
	}
	// Store the URL
	s.data[alias] = url

	// Check for file storage flag, if present - save the request result to a file
	if s.cfg.FileStorage != "" {
		err := Save(s.cfg.FileStorage, alias, url)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// Save writes a URL record to the specified file.
func Save(file string, alias string, url string) error {
	producer, err := NewProducer(file)
	if err != nil {
		return err
	}
	defer producer.Close()

	uuid := generatestring.GenerateUUID()
	event := NewFileStorage(uuid, alias, url)

	if err = producer.WriteEvent(event); err != nil {
		return err
	}
	return nil
}

// Load reads URL records from the specified file and loads them into memory.
func Load(filename string) error {
	d := &Data{
		data: make(map[string]string),
	}

	consumer, err := NewConsumer(filename)
	if err != nil {
		return err
	}
	defer consumer.Close()

	for {
		item, err := consumer.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			return err
		}
		d.data[item.Alias] = item.URL
	}
	return nil
}

// GetStats is not implemented in memory storage.
func (s *Data) GetStats(ctx context.Context) ([]byte, error) {
	readyStats, err := json.Marshal(stats{
		URLs:  len(s.data),
		Users: 0,
	})
	if err != nil {
		return nil, err
	}
	return readyStats, nil
}
