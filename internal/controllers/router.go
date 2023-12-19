package controllers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/middleware/gzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/logger"
)

//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/shortenerURL/internal/controllers UseCase
type UseCase interface {
	DoGet(ctx context.Context, alias string) (string, bool, error)
	DoGetAll(ctx context.Context, userID int, url string) ([]byte, error)
	DoPut(ctx context.Context, url string, uuid int) (string, error)
	DoDel(ctx context.Context, id int, aliases []string) error
	DoHealthcheck() (bool, error)
}
type Controller struct {
	uc  UseCase
	log *zap.Logger
	cfg config.HTTPServer
}

func New(uc UseCase, log *zap.Logger, cfg config.HTTPServer) *Controller {
	return &Controller{uc: uc, log: log, cfg: cfg}
}

func (c *Controller) Router(handler *chi.Mux) *chi.Mux {
	handler.Use(middleware.RequestID)
	handler.Use(mwLogger.New(c.log, c.cfg))
	handler.Use(middleware.Logger)
	handler.Use(gzip.New())

	h := New(c.uc, c.log, c.cfg)

	// Настройка маршрутов с использованием middleware
	handler.Group(func(r chi.Router) {
		r.Get("/{id}", h.Get)
		r.Get("/api/user/urls", h.GetAll)
		r.Get("/ping", h.HealthCheck)
		r.Post("/api/shorten", h.Shorten)
		r.Post("/api/shorten/batch", h.Batch)
		r.Post("/", h.Save)
		r.Delete("/api/user/urls", h.Del)
	})
	return handler
}
