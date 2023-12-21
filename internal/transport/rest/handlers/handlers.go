package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
)

type Handlers struct {
	Get     http.HandlerFunc
	GetAll  http.HandlerFunc
	Shorten http.HandlerFunc
	Save    http.HandlerFunc
	Ping    http.HandlerFunc
	Batch   http.HandlerFunc
	Del     http.HandlerFunc
}

// New создает экземпляр Handlers, инициализируя каждый хендлер
func New(db app.Storage, log *zap.Logger, cfg config.Args) *Handlers {
	if db == nil {
		db, _ = dbstorage.New(cfg.DSN, log)
	}
	return &Handlers{
		Get:     NewGetHandler(db, log, cfg).ServeHTTP,
		GetAll:  NewGetAllHandler(db, log, cfg).ServeHTTP,
		Shorten: NewShortenHandlers(db, log, cfg).ServeHTTP,
		Save:    NewSaveHandlers(db, log, cfg).ServeHTTP,
		Ping:    NewHealtCheck(db).ServeHTTP,
		Batch:   NewBatchHandler(db, log, cfg).ServeHTTP,
		Del:     NewDelURL(db, log, cfg).ServeHTTP,
	}
}
