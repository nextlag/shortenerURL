package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	SaveURL(string, string) error
	Get(string) (string, error)
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

// SaveURL сохраняет значение по ключу
func (s *InMemoryStorage) SaveURL(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Проверка на пустое значение ключа
	if len(key) == 0 {
		return fmt.Errorf("key '%s' cannot be empty", key)
	}
	s.data[key] = value
	return nil
}
