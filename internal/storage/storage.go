package storage

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
	"github.com/nextlag/shortenerURL/internal/lib/storagefile"
	"io"
	"log"
	"sync"
)

// Storage представляет интерфейс для хранилища данных
type (
	Storage interface {
		Get(string) (string, error)
		Put(string, string) error
	}
)

// InMemoryStorage представляет реализацию интерфейса Storage
type InMemoryStorage struct {
	Data  map[string]string
	Mutex sync.Mutex // Мьютекс для синхронизации доступа к данным
}

// NewInMemoryStorage - конструктор для создания нового экземпляра InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Data: make(map[string]string),
	}
}

// Get возвращает значение по ключу
func (s *InMemoryStorage) Get(key string) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	value, ok := s.Data[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return value, nil
}

// Put сохраняет значение по ключу
func (s *InMemoryStorage) Put(key, value string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	// Проверка на пустое значение ключа
	if len(key) == 0 {
		return fmt.Errorf("key '%s' cannot be empty", key)
	}

	// Проверка уникальности данных
	for existingKey, existingValue := range s.Data {
		if existingKey == key || existingValue == value {
			return fmt.Errorf("alias '%s' or URL '%s' already exists", key, value)
		}
	}

	s.Data[key] = value
	err := Save(config.Args.FileStorage, key, value)
	if err != nil {
		return err
	}
	return nil
}

func Save(file string, alias string, originalURL string) error {
	Producer, err := storagefile.NewProducer(file)
	if err != nil {
		log.Fatal(err)
	}
	defer Producer.Close()

	event := &storagefile.Event{
		UUID:  generatestring.GenerateUUID(),
		Alias: alias,
		URL:   originalURL,
	}
	if err := Producer.WriteEvent(event); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *InMemoryStorage) Load(filename string) error {
	Consumer, err := storagefile.NewConsumer(filename)
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
		s.Data[item.Alias] = item.URL
	}
	return nil
}
