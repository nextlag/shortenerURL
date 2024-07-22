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

const fileDel = "del.json"

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
	mutex sync.RWMutex
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
	s.mutex.RLock()
	defer s.mutex.RUnlock()

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
func (s *Data) GetAll(_ context.Context, _ int, _ string) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

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
			err := save(s.cfg.FileStorage, alias, "", userID, s.del[alias].IsDeleted)
			if err != nil {
				return err
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
		err := save(s.cfg.FileStorage, alias, url, userID, s.del[alias].IsDeleted)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// save writes a URL record to the specified file.
func save(file, alias, url string, userID int, isDeleted bool) error {
	uuid := strconv.Itoa(userID)
	if url != "" {
		producer, err := NewProducer(file)
		if err != nil {
			return err
		}
		defer producer.Close()
		event := NewFileStorage(uuid, alias, url)
		if err = producer.WriteEvent(event); err != nil {
			return err
		}
	}
	producerDel, err := NewProducer(fileDel)
	if err != nil {
		return err
	}
	defer producerDel.Close()
	eventDel := NewIsDeleted(uuid, alias, isDeleted)
	if err := producerDel.WriteEventDel(eventDel); err != nil {
		return err
	}
	return nil
}

// Load reads URL records from the specified file and loads them into memory.
func Load(filename string, db *Data) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		consumerDel, err := NewConsumer(fileDel)
		if err != nil {
			errChan <- err
			return
		}
		defer consumerDel.Close()

		delMap := make(map[string]bool)

		for {
			itemDel, err := consumerDel.ReadEventDel()
			if err != nil {
				if err == io.EOF {
					break
				}
				errChan <- err
				return
			}
			if itemDel.StatusDel {
				delMap[itemDel.Alias] = true
			}
		}

		db.mutex.Lock()
		for alias := range delMap {
			db.del[alias] = &dataDel{
				IsDeleted: true,
			}
		}
		db.mutex.Unlock()
	}()

	go func() {
		defer wg.Done()
		consumer, err := NewConsumer(filename)
		if err != nil {
			errChan <- err
			return
		}
		defer consumer.Close()

		for {
			item, err := consumer.ReadEvent()
			if err != nil {
				if err == io.EOF {
					break
				}
				errChan <- err
				return
			}

			db.mutex.Lock()
			if _, exists := db.del[item.Alias]; !exists {
				db.data[item.Alias] = item.URL
			}
			db.mutex.Unlock()
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// GetStats is not implemented in memory storage.
func (s *Data) GetStats(_ context.Context) ([]byte, error) {
	readyStats, err := json.Marshal(stats{
		URLs:  len(s.data),
		Users: 0,
	})
	if err != nil {
		return nil, err
	}
	return readyStats, nil
}
