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
	handler.Use(middleware.Recoverer)

	// Настройка маршрутов с использованием middleware
	handler.Group(func(r chi.Router) {
		r.Get("/{id}", c.Get)
		r.Get("/api/user/urls", c.GetAll)
		r.Get("/ping", c.HealthCheck)
		r.Post("/api/shorten", c.Shorten)
		r.Post("/api/shorten/batch", c.Batch)
		r.Post("/", c.Save)
		r.Delete("/api/user/urls", c.Del)
	})
	return handler
}
