package storage

type Storage interface {
	Get(string) (string, bool)
	Put(string, string)
}

type Database map[string]string

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

func (s *InMemoryStorage) Put(key, value string) {
	s.data[key] = value
}
