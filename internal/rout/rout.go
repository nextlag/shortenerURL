package rout

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nextlag/shortenerURL/internal/handlers/httpserver"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/zapLogger"
	"github.com/nextlag/shortenerURL/internal/storage"
	"go.uber.org/zap"
)

func SetupRouter(db storage.Storage, log *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // добавляем уникальный идентификатор
	router.Use(middleware.Logger)    // добавляем вывод стандартного логгера

	// Создание экземпляра middleware.Logger
	mw := mwLogger.NewLogger(log)

	// Настройка маршрутов с использованием middleware
	router.With(mw).Get("/{id}", httpserver.GetHandler(db))
	router.With(mw).Post("/api/shorten", httpserver.Shorten(log, db))
	router.With(mw).Post("/", httpserver.Save(db))

	return router
}
