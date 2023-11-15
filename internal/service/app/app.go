package app

import (
	"flag"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// TODO:

// Storage представляет интерфейс для хранилища данных
type Storage interface {
	Get(string) (string, error)
	Put(string, string) error
	Load(string) error
}

type App struct {
	Stor Storage           `json:"stor,omitempty"`
	Log  *zap.Logger       `json:"log,omitempty"`
	Cfg  config.ConfigHTTP `json:"cfg,omitempty"`
}

func New() *App {
	flag.Parse()
	return &App{
		Stor: storage.New(),
		Log:  lg.New(),
		Cfg:  config.Config,
	}
}

type DBStorage interface {
	Stop() error
	CheckConnection() bool
	CreateTable() error
}
