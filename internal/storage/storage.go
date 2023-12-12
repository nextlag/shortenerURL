package storage

import (
	"context"
	"fmt"
	"io"
	"sync"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/filestorage"
	"github.com/nextlag/shortenerURL/internal/utils/generatestring"
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

func (s *Data) Get(_ context.Context, alias string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Получение из файла или памяти
	url, ok := s.data[alias]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", alias)
	}
	return url, nil
}

func (s *Data) GetAll(_ context.Context, _ int, _ string) ([]byte, error) {
	return []byte("Memory storage can't operate with user IDs"), nil
}

func (s *Data) CheckConnection() bool {
	return true
}

// Put сохраняет значение по ключу
func (s *Data) Put(_ context.Context, url string, _ int) (string, error) {
	var log = app.New().Log
	var cfg = app.New().Cfg
	s.mutex.Lock()
	defer s.mutex.Unlock()

	alias := generatestring.NewRandomString(8)

	// Проверяем существование ключа
	if _, exists := s.data[alias]; exists {
		return "", fmt.Errorf("alias '%s/%s' already exists", cfg.BaseURL, alias)
	}

	// Проверяем существование значения
	for k, v := range s.data {
		if v == url {
			log.Info("response", zap.String("ulr", v), zap.String("alias", k))
			return k, nil
		}
	}
	// Запись url
	s.data[alias] = url

	// Проверка на существование флага -f, если есть - сохранить результат запроса в файл
	if cfg.FileStorage != "" {
		err := Save(log, cfg.FileStorage, alias, url)
		if err != nil {
			return alias, err
		}
	}
	return alias, nil
}

func Save(log *zap.Logger, file string, alias string, url string) error {
	Producer, err := filestorage.NewProducer(file)
	if err != nil {
		return err
	}
	defer Producer.Close()

	uuid := generatestring.GenerateUUID()
	event := filestorage.New(uuid, alias, url)

	if err := Producer.WriteEvent(event); err != nil {
		return err
	}

	log.Info("Data.Put", zap.Any("Save", event))

	return nil
}

func Load(filename string) (*Data, error) {
	s := &Data{
		data: make(map[string]string),
	}

	Consumer, err := filestorage.NewConsumer(filename)
	if err != nil {
		return nil, err
	}

	defer Consumer.Close()

	for {
		item, err := Consumer.ReadEvent()
		if err != nil {
			if err == io.EOF {
				break // Достигнут конец файла
			}
			return nil, err
		}
		s.data[item.Alias] = item.URL
	}
	return s, nil
}
