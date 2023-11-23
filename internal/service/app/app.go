package app

import (
	"flag"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/storage/database/dbstorage"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Storage представляет интерфейс для хранилища данных
type Storage interface {
	Get(string) (string, error)
	Put(string, string) error
}

type App struct {
	Stor Storage
	db   dbstorage.DBStorage
	Log  *zap.Logger
	Cfg  config.ConfigHTTP
}

func New() *App {
	flag.Parse()
	return &App{
		Stor: storage.New(),
		Log:  lg.New(),
		Cfg:  config.Config,
	}
}
