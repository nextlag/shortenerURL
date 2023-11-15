package storage

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/database/dbstorage"
	"github.com/nextlag/shortenerURL/internal/storage/filestorage"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Data представляет реализацию интерфейса Storage
type Data struct {
	data  map[string]string
	mutex sync.Mutex // Мьютекс для синхронизации доступа к данным
}

// New - конструктор для создания нового экземпляра Data
func New() *Data {
	return &Data{
		data: make(map[string]string),
	}
}

func (s *Data) Get(alias string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	logger := lg.New()

	if config.Config.DSN != "" {
		// Получение из базы данных
		db, err := dbstorage.New(config.Config.DSN)
		if err != nil {
			return "", fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.Stop()

		url, err := db.Get(alias)
		if err != nil {
			logger.Error("failed to get data from DB", zap.Error(err))
			return "", err
		}
		return url.URL, nil
	}

	// Получение из файла или памяти
	url, ok := s.data[alias]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", alias)
	}
	return url, nil
}

// Put сохраняет значение по ключу
func (s *Data) Put(alias, url string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if config.Config.DSN != "" {
		// Сохранение в базу данных
		db, err := dbstorage.New(config.Config.DSN)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.Stop()

		shortURL := &dbstorage.ShortURL{
			URL:       url,
			Alias:     alias,
			CreatedAt: time.Now(),
		}

		_, err = db.Insert(shortURL)
		if err != nil {
			return fmt.Errorf("failed to insert short URL into database: %w", err)
		}
		return nil
	}

	// Проверка на пустое значение ключа
	if len(alias) == 0 {
		return fmt.Errorf("key '%s' cannot be empty", alias)
	}
	// Проверка уникальности данных
	for existingKey, existingValue := range s.data {
		if existingKey == alias || existingValue == url {
			return fmt.Errorf("alias '%s' or URL '%s' already exists", alias, url)
		}
	}
	s.data[alias] = url
	err := Save(config.Config.FileStorage, alias, url)
	if err != nil {
		return err
	}
	return nil
}

func Save(file string, alias string, url string) error {
	Producer, err := filestorage.NewProducer(file)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()
	uuid := generatestring.GenerateUUID()
	event := filestorage.New(uuid, alias, url)
	if err := Producer.WriteEvent(event); err != nil {
		log.Fatal(err)
	}
	logger := lg.New()
	logger.Info("add_request", zap.Any("data", event))
	return nil
}

func (s *Data) Load(filename string) error {
	Consumer, err := filestorage.NewConsumer(filename)
	if err != nil {
		return err
	}
	defer Consumer.Close()

	for {
		item, err := Consumer.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break // Достигнут конец файла
			}
			return err
		}
		s.data[item.Alias] = item.URL
	}
	return nil
}
