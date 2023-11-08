package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/service/app"
)

type Handlers struct {
	Get     http.HandlerFunc
	Shorten http.HandlerFunc
	Save    http.HandlerFunc
}

// New создает экземпляр Handlers, инициализируя каждый хендлер
func New(log *zap.Logger, db app.Storage) *Handlers {
	return &Handlers{
		Get:     GetHandler(db),
		Shorten: Shorten(log, db),
		Save:    Save(db),
	}
}
