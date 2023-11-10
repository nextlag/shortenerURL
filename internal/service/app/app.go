package app

import (
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
	Stor Storage
	Log  *zap.Logger
	Cfg  config.ArgsHTTP
}

func New() *App {
	return &App{
		Stor: storage.New(),
		Log:  lg.New(),
		Cfg:  config.Args,
	}
}
