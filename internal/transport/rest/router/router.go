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
	// router.Use(middleware.Logger)    // добавляем вывод стандартного логгера

	// Создание экземпляра middleware.Logger
	mw := mwLogger.New(log)
	h := handlers.New(log, db)
	// Настройка маршрутов с использованием middleware
	router.With(mw).Get("/{id}", h.Get)
	router.With(mw).Post("/api/shorten", h.Shorten)
	router.With(mw).Get("/ping", h.Ping)
	router.With(mw).Post("/", h.Save)

	return router
}
