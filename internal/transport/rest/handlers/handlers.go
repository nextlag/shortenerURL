package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
)

type Handler struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

type Handlers struct {
	Get     http.HandlerFunc
	GetAll  http.HandlerFunc
	Shorten http.HandlerFunc
	Save    http.HandlerFunc
	Ping    http.HandlerFunc
	Batch   http.HandlerFunc
}

// New создает экземпляр Handlers, инициализируя каждый хендлер
func New(log *zap.Logger, stor app.Storage, db *dbstorage.DBStorage) *Handlers {
	if db == nil {
		db, _ = dbstorage.New(config.Config.DSN)
	}
	return &Handlers{
		Get:     GetHandler(stor),
		GetAll:  NewGetAllHandler(log, stor).ServeHTTP,
		Shorten: Shorten(log, stor),
		Save:    NewSaveHandlers(db).ServeHTTP,
		Ping:    NewHealCheck(db).ServeHTTP,
		Batch:   NewBatchHandler(log, stor).ServeHTTP,
	}
}
