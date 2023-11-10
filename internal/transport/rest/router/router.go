package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/transport/rest/handlers"
	mwLogger "github.com/nextlag/shortenerURL/internal/transport/rest/middleware/zaplogger"
)

func SetupRouter(db app.Storage, log *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // добавляем уникальный идентификатор
	router.Use(mwLogger.New(log))    // добавляем middleware.Logger

	h := handlers.New(log, db)

	// Настройка маршрутов с использованием middleware
	router.With(mwLogger.New(log)).Get("/{id}", h.Get)
	router.With(mwLogger.New(log)).Post("/api/shorten", h.Shorten)
	router.With(mwLogger.New(log)).Get("/ping", h.Ping)
	router.With(mwLogger.New(log)).Post("/", h.Save)

	return router
}
