package storage

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/lib/generatestring"
	"github.com/nextlag/shortenerURL/internal/lib/storagefile"
	"log"
	"sync"
)

// Storage представляет интерфейс для хранилища данных
type Storage interface {
	Get(string) (string, error)
	Put(string, string) error
	Save(string, string, string) error
}

// InMemoryStorage представляет реализацию интерфейса Storage
type InMemoryStorage struct {
	data  map[string]string
	mutex sync.Mutex // Мьютекс для синхронизации доступа к данным
}

// NewInMemoryStorage - конструктор для создания нового экземпляра InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

// Get возвращает значение по ключу
func (s *InMemoryStorage) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return value, nil
}

// Put сохраняет значение по ключу
func (s *InMemoryStorage) Put(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Проверка на пустое значение ключа
	if len(key) == 0 {
		return fmt.Errorf("key '%s' cannot be empty", key)
	}
	s.data[key] = value
	return nil
}

func (s *InMemoryStorage) Save(file string, alias string, originalURL string) error {
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
