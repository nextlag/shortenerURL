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
func New(log *zap.Logger, db app.Storage) *Handlers {
	if db == nil {
		db, _ = dbstorage.New(config.Config.DSN)
	}
	return &Handlers{
		Get:     GetHandler(db),
		GetAll:  NewGetAllHandler(log, db).ServeHTTP,
		Shorten: Shorten(log, db),
		Save:    NewSaveHandlers(db).ServeHTTP,
		Ping:    NewHealCheck(db).ServeHTTP,
		Batch:   NewBatchHandler(db).ServeHTTP,
	}
}
