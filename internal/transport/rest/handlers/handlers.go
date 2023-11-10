package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	dbstorage "github.com/nextlag/shortenerURL/internal/database/psql"
	"github.com/nextlag/shortenerURL/internal/service/app"
)

type Handlers struct {
	Get     http.HandlerFunc
	Shorten http.HandlerFunc
	Save    http.HandlerFunc
	Ping    http.HandlerFunc
}

// New создает экземпляр Handlers, инициализируя каждый хендлер
func New(log *zap.Logger, stor app.Storage, db *dbstorage.DBStorage) *Handlers {
	if db == nil {
		db, _ = dbstorage.New(config.Args.Psql)
	}
	pingHandler := NewHealCheck(db)
	return &Handlers{
		Get:     GetHandler(stor),
		Shorten: Shorten(log, stor),
		Save:    Save(stor),
		Ping:    pingHandler.healCheck,
	}
}
