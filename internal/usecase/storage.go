package usecase

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
)

type Data struct {
	data  map[string]string
	mutex sync.Mutex // Мьютекс для синхронизации доступа к данным
}

// New - конструктор для создания нового экземпляра Data
func NewData() *Data {
	return &Data{
		data: make(map[string]string),
	}
}

func (s *Data) Get(_ context.Context, alias string) (string, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Получение из файла или памяти
	url, ok := s.data[alias]
	if !ok {
		return "", false, fmt.Errorf("key '%s' not found", alias)
	}
	return url, false, nil
}

func (s *Data) GetAll(context.Context, int, string) ([]byte, error) {
	return []byte("Memory storage can't operate with user IDs"), nil
}

func (s *Data) Healthcheck() (bool, error) {
	return true, nil
}

func (s *Data) Del(_ context.Context, _ int, _ []string) error {
	return nil
}

// Put сохраняет значение по ключу
func (s *Data) Put(_ context.Context, url string, _ int) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	cfg := config.Cfg
	alias := generatestring.NewRandomString(8)

	// Проверяем существование ключа
	if _, exists := s.data[alias]; exists {
		return "", fmt.Errorf("alias '%s/%s' already exists", cfg.BaseURL, alias)
	}

	// Проверяем существование значения
	for k, v := range s.data {
		if v == url {
			return k, nil
		}
	}
	// Запись url
	s.data[alias] = url

	// Проверка на существование флага -f, если есть - сохранить результат запроса в файл
	if cfg.FileStorage != "" {
		err := Save(cfg.FileStorage, alias, url)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

func Save(file string, alias string, url string) error {
	Producer, err := NewProducer(file)
	if err != nil {
		return err
	}
	defer Producer.Close()

	uuid := generatestring.GenerateUUID()
	event := NewFileStorage(uuid, alias, url)

	if err := Producer.WriteEvent(event); err != nil {
		return err
	}
	return nil
}

func Load(filename string) error {
	d := &Data{
		data: make(map[string]string),
	}

	Consumer, err := NewConsumer(filename)
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
		d.data[item.Alias] = item.URL
	}
	return nil
}
