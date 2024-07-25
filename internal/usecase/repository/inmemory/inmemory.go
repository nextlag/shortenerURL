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
	URL       string
	IsDeleted bool
}

// Data represents the in-memory data storage structure.
type Data struct {
	data  map[string]*dataDel
	log   *zap.Logger
	cfg   *configuration.Config
	mutex sync.RWMutex
}

// New creates a new instance of Data.
func New(cfg *configuration.Config, log *zap.Logger) (*Data, error) {
	return &Data{
		data: make(map[string]*dataDel),
		log:  log,
		cfg:  cfg,
	}, nil
}

// Get retrieves a URL by its alias from the in-memory storage.
func (s *Data) Get(_ context.Context, alias string) (*entity.URL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	delInfo, ok := s.data[alias]
	if !ok || delInfo.IsDeleted {
		return nil, fmt.Errorf("key '%s' not found", alias)
	}

	return &entity.URL{
		Alias: alias,
		URL:   delInfo.URL,
	}, nil
}

// GetAll retrieves all URLs for a given user.
func (s *Data) GetAll(_ context.Context, _ int, host string) ([]*entity.URL, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var userUrls []*entity.URL
	for alias, delInfo := range s.data {
		if !delInfo.IsDeleted {
			userUrls = append(userUrls, &entity.URL{
				Alias: fmt.Sprintf("%s/%s", host, alias),
				URL:   delInfo.URL,
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
		if delInfo, exists := s.data[alias]; exists {
			delInfo.IsDeleted = true
			err := save(s.cfg.FileStorage, alias, delInfo.URL, userID, delInfo.IsDeleted)
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
		if v.URL == url {
			return k, nil
		}
	}

	s.data[alias] = &dataDel{
		UserID:    strconv.Itoa(userID),
		URL:       url,
		IsDeleted: false,
	}

	if s.cfg.FileStorage != "" {
		err := save(s.cfg.FileStorage, alias, url, userID, s.data[alias].IsDeleted)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// save writes URL record and deletion status to the specified files.
func save(file, alias, url string, userID int, isDeleted bool) error {
	uuid := strconv.Itoa(userID)
	if url != "" {
		producer, err := NewProducer(file)
		if err != nil {
			return err
		}
		defer producer.Close()
		event := NewFileStorage(uuid, alias, url)
		if err = WriteEvent(producer, event); err != nil {
			return err
		}
	}

	producerDel, err := NewProducer(fileDel)
	if err != nil {
		return err
	}
	defer producerDel.Close()
	eventDel := NewIsDeleted(uuid, alias, isDeleted)
	if err = WriteEvent(producerDel, eventDel); err != nil {
		return err
	}
	return nil
}

// Load reads URL records and deletion statuses from the specified files and loads them into memory.
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
			itemDel, err := ReadEvent[IsDeleted](consumerDel)
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
		defer db.mutex.Unlock()

		for alias := range delMap {
			if delInfo, exists := db.data[alias]; exists {
				delInfo.IsDeleted = true
			} else {
				db.data[alias] = &dataDel{IsDeleted: true}
			}
		}
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
			item, err := ReadEvent[FileStorage](consumer)
			if err != nil {
				if err == io.EOF {
					break
				}
				errChan <- err
				return
			}

			db.mutex.Lock()
			if delInfo, exists := db.data[item.Alias]; !exists || !delInfo.IsDeleted {
				db.data[item.Alias] = &dataDel{
					UserID: item.UUID,
					URL:    item.URL,
				}
			}
			db.mutex.Unlock()
		}
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// GetStats retrieves statistics on the number of URLs and users.
func (s *Data) GetStats(_ context.Context) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	urlCount := 0
	userMap := make(map[string]struct{})

	for _, delInfo := range s.data {
		if !delInfo.IsDeleted {
			urlCount++
			userMap[delInfo.UserID] = struct{}{}
		}
	}
	userCount := len(userMap)

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
