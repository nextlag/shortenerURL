package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/database/dbstorage"
)

type Handlers struct {
	Get     http.HandlerFunc
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
	pingHandler := NewHealCheck(db)
	return &Handlers{
		Get:     GetHandler(stor),
		Shorten: Shorten(log, stor),
		Save:    Save(stor),
		Ping:    pingHandler.healCheck,
		Batch:   NewBatchHandler(log, db).ServeHTTP,
	}
}

// BatchHandler представляет хендлер для сокращения нескольких URL.
type BatchHandler struct {
	log *zap.Logger
	db  *dbstorage.DBStorage
}

// NewBatchHandler создает новый экземпляр BatchHandler.
func NewBatchHandler(log *zap.Logger, db *dbstorage.DBStorage) *BatchHandler {
	return &BatchHandler{
		log: log,
		db:  db,
	}
}
