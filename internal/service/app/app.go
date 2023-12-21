package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
)

// Storage представляет интерфейс для хранилища данных
//
//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/shortenerURL/internal/service/app Storage
type Storage interface {
	Get(context.Context, string) (string, bool, error)
	GetAll(ctx context.Context, userID int, url string) ([]byte, error)
	Put(context.Context, string, int) (string, error)
	Del(context.Context, int, []string) error
	Healtcheck() bool
}

type App struct {
	DB  Storage
	Log *zap.Logger
	Cfg config.Args
}

func New() *App {
	return &App{
		Log: lg.New(),
		Cfg: config.Config,
	}
}
