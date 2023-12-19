package storage

import (
	"encoding/json"
	"os"
)

type FileStorage struct {
	UUID  string `json:"uuid"`                        // UUID, генерация uuid
	Alias string `json:"alias,omitempty"`             // Alias, пользовательский псевдоним для короткой ссылки (необязательный).
	URL   string `json:"url" validate:"required,url"` // URL, который нужно сократить, должен быть валидным URL.

}

func NewFileStorage(uuid, alias, url string) *FileStorage {
	return &FileStorage{
		UUID:  uuid,
		Alias: alias,
		URL:   url,
	}
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *FileStorage) error {
	return p.encoder.Encode(event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*FileStorage, error) {
	event := &FileStorage{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
