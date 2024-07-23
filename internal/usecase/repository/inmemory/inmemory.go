package inmemory

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
	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/usecase/repository/models"
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

// New creates a new instance of Data.
func New(cfg *configuration.Config, log *zap.Logger) (*Data, error) {
	return &Data{
		data: make(map[string]string),
		del:  make(map[string]*dataDel),
		log:  log,
		cfg:  cfg,
	}, nil
}

// Get retrieves a URL by its alias from the in-memory storage.
func (s *Data) Get(_ context.Context, alias string) (*entity.URL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	url, ok := s.data[alias]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", alias)
	}

	if delInfo, exists := s.del[alias]; exists && delInfo.IsDeleted {
		return nil, fmt.Errorf("key '%s' is deleted", alias)
	}

	// Since we don't store full entity.URL details in memory, create a placeholder
	return &entity.URL{
		Alias: alias,
		URL:   url,
	}, nil
}

// GetAll retrieves all URLs for a given user.
func (s *Data) GetAll(_ context.Context, _ int, host string) ([]*entity.URL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var userUrls []*entity.URL
	for alias, url := range s.data {
		if delInfo, exists := s.del[alias]; !exists || !delInfo.IsDeleted {
			userUrls = append(userUrls, &entity.URL{
				Alias: fmt.Sprintf("%s/%s", host, alias),
				URL:   url,
			})
		}
	}
	return userUrls, nil
}

// Healthcheck checks if a file exists at the specified path and is accessible for reading and writing.
func (s *Data) Healthcheck() (bool, error) {
	filePath := s.cfg.FileStorage
	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			file, createErr := os.Create(filePath)
			if createErr != nil {
				s.log.Error("Healthcheck: unable to create file", zap.Error(createErr))
				return false, errors.New("file does not exist and cannot be created")
			}
			file.Close()
			s.log.Info("Healthcheck: file created successfully", zap.String("file", filePath))
			return true, nil
		}
		s.log.Error("Healthcheck: error stating file", zap.Error(err))
		return false, err
	}

	file, openErr := os.OpenFile(filePath, os.O_RDWR, 0666)
	if openErr != nil {
		s.log.Error("Healthcheck: unable to open file", zap.Error(openErr))
		return false, openErr
	}
	file.Close()

	s.log.Info("Healthcheck: file is accessible", zap.String("file", filePath))
	return true, nil
}

// Del marks URLs as deleted in the in-memory storage.
func (s *Data) Del(_ context.Context, userID int, aliases []string) error {
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

// GetStats retrieves statistics on the number of URLs and users.
func (s *Data) GetStats(_ context.Context) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	urlCount := len(s.data)
	userCount := 0
	userMap := make(map[string]struct{})

	for _, delInfo := range s.del {
		if !delInfo.IsDeleted {
			userMap[delInfo.UserID] = struct{}{}
		}
	}
	userCount = len(userMap)

	stats := models.Stats{
		URLs:  urlCount,
		Users: userCount,
	}

	readyStats, err := json.Marshal(stats)
	if err != nil {
		return nil, err
	}
	return readyStats, nil
}
