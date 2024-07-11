package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/pkg/tools/generatestring"
)

type dataDel struct {
	UserID    string
	Alias     string
	IsDeleted bool
}

// Data represents the in-memory data storage structure.
type Data struct {
	data  map[string]string
	del   map[string]*dataDel
	log   *zap.Logger
	cfg   *configuration.Config
	mutex sync.Mutex
}

// NewData creates a new instance of Data.
func NewData(log *zap.Logger, cfg *configuration.Config) *Data {
	return &Data{
		data: make(map[string]string),
		del:  make(map[string]*dataDel),
		log:  log,
		cfg:  cfg,
	}
}

// Get retrieves a URL by its alias from the in-memory storage.
func (s *Data) Get(_ context.Context, alias string) (string, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	url, ok := s.data[alias]
	if !ok {
		return "", false, fmt.Errorf("key '%s' not found", alias)
	}

	if delInfo, exists := s.del[alias]; exists && delInfo.IsDeleted {
		return "", false, fmt.Errorf("key '%s' is deleted", alias)
	}

	return url, false, nil
}

// GetAll retrieves all URLs for a given user.
func (s *Data) GetAll(_ context.Context, userID int, _ string) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	userUrls := make(map[string]string)
	for alias, urlData := range s.data {
		if delInfo, exists := s.del[alias]; !exists || !delInfo.IsDeleted {
			userUrls[alias] = urlData
		}
	}
	return json.Marshal(userUrls)
}

// Healthcheck checks if a file exists at the specified path.
func (s *Data) Healthcheck() (bool, error) {
	filePath := s.cfg.FileStorage
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, errors.New("file does not exist")
		}
		return false, err
	}
	return true, nil
}

// Del marks URLs as deleted in the in-memory storage.
func (s *Data) Del(userID int, aliases []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, alias := range aliases {
		if _, exists := s.data[alias]; exists {
			s.del[alias] = &dataDel{
				UserID:    strconv.Itoa(userID),
				Alias:     alias,
				IsDeleted: true,
			}
		}
	}
	return nil
}

// Put saves a URL with a generated alias in the in-memory storage.
func (s *Data) Put(_ context.Context, url string, alias string, userID int) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if alias == "" {
		alias = generatestring.NewRandomString(8)
	}

	if _, exists := s.data[alias]; exists {
		return "", fmt.Errorf("alias '%s/%s' already exists", s.cfg.BaseURL, alias)
	}

	for k, v := range s.data {
		if v == url {
			return k, nil
		}
	}

	s.data[alias] = url

	s.del[alias] = &dataDel{
		UserID:    strconv.Itoa(userID),
		Alias:     alias,
		IsDeleted: false,
	}

	if s.cfg.FileStorage != "" {
		err := save(s.cfg.FileStorage, alias, url, userID)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// save writes a URL record to the specified file.
func save(file string, alias string, url string, userID int) error {
	producer, err := NewProducer(file)
	if err != nil {
		return err
	}
	defer producer.Close()

	uuid := strconv.Itoa(userID)
	event := NewFileStorage(uuid, alias, url)

	if err = producer.WriteEvent(event); err != nil {
		return err
	}
	return nil
}

// Load reads URL records from the specified file and loads them into memory.
func Load(filename string, db *Data) error {
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
		db.data[item.Alias] = item.URL
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
