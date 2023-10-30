package storage

import (
	"fmt"
	"io"
	"log"
	"sync"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
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

// Get возвращает значение по ключу
func (s *Data) Get(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	value, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return value, nil
}

// Put сохраняет значение по ключу
func (s *Data) Put(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Проверка на пустое значение ключа
	if len(key) == 0 {
		return fmt.Errorf("key '%s' cannot be empty", key)
	}

	// Проверка уникальности данных
	for existingKey, existingValue := range s.data {
		if existingKey == key || existingValue == value {
			return fmt.Errorf("alias '%s' or URL '%s' already exists", key, value)
		}
	}

	s.data[key] = value
	err := Save(config.Args.FileStorage, key, value)
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
