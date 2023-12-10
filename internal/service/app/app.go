package app

import (
	"context"
	"flag"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Storage представляет интерфейс для хранилища данных
//
//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/shortenerURL/internal/service/app Storage
type Storage interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, int) (string, error)
	GetAll(ctx context.Context, id int, url string) ([]byte, error)
}

type App struct {
	Stor Storage
	Log  *zap.Logger
	Cfg  config.Args
}

func New() *App {
	flag.Parse()
	return &App{
		Stor: storage.New(),
		Log:  lg.New(),
		Cfg:  config.Config,
	}
}
