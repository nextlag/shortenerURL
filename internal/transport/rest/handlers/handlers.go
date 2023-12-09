package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
)

type Handlers struct {
	Get        http.HandlerFunc
	GetAllURLs http.HandlerFunc
	Shorten    http.HandlerFunc
	Save       http.HandlerFunc
	Ping       http.HandlerFunc
	Batch      http.HandlerFunc
}

// New создает экземпляр Handlers, инициализируя каждый хендлер
func New(log *zap.Logger, stor app.Storage, db *dbstorage.DBStorage) *Handlers {
	if db == nil {
		db, _ = dbstorage.New(config.Config.DSN)
	}
	pingHandler := NewHealCheck(db)
	return &Handlers{
		Get:        GetHandler(stor),
		GetAllURLs: NewGetAllHandler(log, stor).ServeHTTP,
		Shorten:    Shorten(log, stor),
		Save:       Save(stor),
		Ping:       pingHandler.healCheck,
		Batch:      NewBatchHandler(log, stor).ServeHTTP,
	}
}
