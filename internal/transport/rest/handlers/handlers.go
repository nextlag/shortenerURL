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
func New(log *zap.Logger, db app.Storage, dbStorage *dbstorage.DBStorage) *Handlers {
	dbStorage, _ = dbstorage.New(config.Args.Psql)
	pingHandler := NewHealCheck(dbStorage)
	return &Handlers{
		Get:     GetHandler(db),
		Shorten: Shorten(log, db),
		Save:    Save(db),
		Ping:    pingHandler.healCheck,
	}
}
