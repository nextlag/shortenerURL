package http

import (
	"context"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/internal/entity"
	"github.com/nextlag/shortenerURL/internal/middleware/gzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/logger"
)

// UseCase defines the interface for the application's use case layer.
//
//go:generate mockgen -destination=mocks/mocks.go -package=mocks github.com/nextlag/shortenerURL/internal/controllers/http UseCase
type UseCase interface {
	DoGet(ctx context.Context, alias string) (*entity.URL, error)
	DoGetAll(ctx context.Context, userID int, url string) ([]*entity.URL, error)
	DoPut(ctx context.Context, url string, alias string, uuid int) (string, error)
	DoDel(ctx context.Context, id int, aliases []string)
	DoHealthcheck() (bool, error)
	DoGetStats(ctx context.Context) ([]byte, error)
}

// Controller represents the application's HTTP controller.
type Controller struct {
	uc  UseCase
	wg  *sync.WaitGroup
	log *zap.Logger
	cfg *configuration.Config
}

// New creates a new Controller.
func New(uc UseCase, wg *sync.WaitGroup, cfg *configuration.Config, log *zap.Logger) *Controller {
	return &Controller{uc: uc, wg: wg, cfg: cfg, log: log}
}

// Controller sets up the application's HTTP routing and middleware.
func (c *Controller) Controller(handler *chi.Mux) *chi.Mux {
	handler.Use(middleware.RequestID)
	handler.Use(mwLogger.New(c.log, c.cfg))
	handler.Use(middleware.Logger)
	handler.Use(gzip.New())
	handler.Use(middleware.Recoverer)

	// Set up routes with middleware
	handler.Group(func(r chi.Router) {
		r.Get("/{id}", c.Get)
		r.Get("/api/user/urls", c.GetAll)
		r.Get("/ping", c.HealthCheck)
		r.Get("/api/internal/stats", c.GetStatsHandler)
		r.Post("/api/shorten", c.Shorten)
		r.Post("/api/shorten/batch", c.Batch)
		r.Post("/", c.Save)
		r.Delete("/api/user/urls", c.Del)
	})

	// Add pprof routes
	handler.Route("/debug/pprof", func(r chi.Router) {
		r.Handle("/", http.HandlerFunc(pprof.Index))
		r.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/trace", http.HandlerFunc(pprof.Trace))
	})

	return handler
}
