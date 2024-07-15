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

// UrlData represents the structure for storing URL data.
type UrlData struct {
	URL       string
	IsDeleted bool
	UserID    int
}

// Data represents the in-memory data storage structure.
type Data struct {
	data  map[string]UrlData
	log   *zap.Logger
	cfg   *configuration.Config
	mutex sync.Mutex // Mutex for synchronizing access to data
}

// NewData creates a new instance of Data.
func NewData(log *zap.Logger, cfg *configuration.Config) *Data {
	return &Data{
		data: make(map[string]UrlData),
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
	return url.URL, false, nil
}

// GetAll retrieves all URLs for a given user.
func (s *Data) GetAll(_ context.Context, userID int, _ string) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	userUrls := make(map[string]string)
	for alias, urlData := range s.data {
		if urlData.UserID == userID && !urlData.IsDeleted {
			userUrls[alias] = urlData.URL
		}
	}
	return json.Marshal(userUrls)
}

// Healthcheck always returns true for in-memory storage.
func (s *Data) Healthcheck() (bool, error) {
	return true, nil
}

// Put saves a URL with a generated alias in the in-memory storage.
func (s *Data) Put(_ context.Context, url string, alias string, userID int) (string, error) {
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
		if v.URL == url {
			return k, nil
		}
	}
	// Store the URL
	s.data[alias] = UrlData{
		URL:       url,
		UserID:    userID,
		IsDeleted: false}

	// Check for file storage flag, if present - save the request result to a file
	if s.cfg.FileStorage != "" {
		err := Save(s.cfg.FileStorage, alias, url, userID)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// Del marks URLs as deleted for a given user.
func (s *Data) Del(userID int, aliases []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, alias := range aliases {
		if urlData, exists := s.data[alias]; exists && urlData.UserID == userID {
			urlData.IsDeleted = true
			s.data[alias] = urlData
		}
	}
	return nil
}

// Save writes a URL record to the specified file.
func Save(file string, alias string, url string, userID int) error {
	producer, err := NewProducer(file)
	if err != nil {
		return err
	}
	defer producer.Close()

	// uuid := generatestring.GenerateUUID()
	event := NewFileStorage(alias, url, userID)

	if err = producer.WriteEvent(event); err != nil {
		return err
	}
	return nil
}

// Load reads URL records from the specified file and loads them into memory.
func Load(filename string) error {
	// Создаем Data без указания логгера и конфигурации
	d := &Data{
		data: make(map[string]UrlData),
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
				break // Достигнут конец файла
			}
			return err
		}

		// Загружаем URL данные в память
		d.data[item.Alias] = UrlData{
			URL:       item.URL,
			IsDeleted: false,
			UserID:    0,
		}
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
