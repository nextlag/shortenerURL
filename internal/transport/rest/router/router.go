package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/transport/rest/handlers"
	mwLogger "github.com/nextlag/shortenerURL/internal/transport/rest/middleware/zaplogger"
)

func SetupRouter(db app.Storage, log *zap.Logger, cfg config.Args) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // добавляем уникальный идентификатор

	h := handlers.New(db, log, cfg)

	// Настройка маршрутов с использованием middleware
	router.With(mwLogger.New(log, cfg)).Route("/", func(r chi.Router) {
		r.Get("/{id}", h.Get)
		r.Get("/api/user/urls", h.GetAll)
		r.Get("/ping", h.Ping)
		r.Post("/api/shorten", h.Shorten)
		r.Post("/api/shorten/batch", h.Batch)
		r.Post("/", h.Save)
		r.Delete("/api/user/urls", h.Del)
	})
	return router
}
