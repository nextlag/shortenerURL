package storage

type Storage interface {
	Get(string) (string, bool)
	Put(string, string)
}

type Database map[string]string

// InMemoryStorage - представляет реализацию интерфейса Storage
type InMemoryStorage struct {
	data map[string]string
}

// NewInMemoryStorage - конструктор для создания нового экземпляра InMemoryStorage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

// Get Метод - пытается получить значение по ключу из data и возвращает его вместе с флагом, указывающим, было ли значение найдено
func (s *InMemoryStorage) Get(key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

// Put - добавляет или обновляет значение по ключу
func (s *InMemoryStorage) Put(key, value string) {
	s.data[key] = value
}
