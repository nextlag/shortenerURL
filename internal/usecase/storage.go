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
			if s.cfg.FileDel != "" {
				err := save(s.cfg.FileDel, s.cfg.FileDel, alias, "", userID, s.del[alias].IsDeleted)
				if err != nil {
					return err
				}
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

	if s.cfg.FileStorage != "" && s.cfg.FileDel != "" {
		err := save(s.cfg.FileStorage, s.cfg.FileDel, alias, url, userID, s.del[alias].IsDeleted)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

// save writes a URL record to the specified file.
func save(file, fileDel, alias, url string, userID int, isDeleted bool) error {
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
func Load(filename, fileDel string, db *Data) error {
	consumer, err := NewConsumer(filename)
	if err != nil {
		return err
	}
	defer consumer.Close()

	consumerDel, err := NewConsumer(fileDel)
	if err != nil {
		return err
	}
	defer consumerDel.Close()

	delMap := make(map[string]bool)

	// Чтение файла удалений и заполнение карты
	for {
		itemDel, err := consumerDel.ReadEventDel()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if itemDel.StatusDel {
			delMap[itemDel.Alias] = true
		}
	}

	// Чтение основного файла и загрузка данных в память
	for {
		item, err := consumer.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			return err
		}

		// Если для alias существует запись с status_del равным true, пропускаем его
		if _, exists := delMap[item.Alias]; exists {
			continue
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
