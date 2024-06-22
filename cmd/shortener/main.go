package main

import (
	"context"
	"errors"
	"flag"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/controllers"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

// setupServer initializes and returns a new HTTP server configured with the provided router.
func setupServer(router http.Handler) *http.Server {
	return &http.Server{
		Addr:    config.Cfg.Host,
		Handler: router,
	}
}

// main is the entry point of the application.
func main() {
	// Load configuration
	if err := config.MakeConfig(); err != nil {
		stdLog.Fatal(err)
	}
	flag.Parse()

	var (
		log = logger.SetupLogger()
		cfg = config.Cfg
		uc  *usecase.UseCase
	)

	// Log initialization details
	log.Info(
		"initialized flags",
		zap.String("-a", cfg.Host),
		zap.String("-b", cfg.BaseURL),
		zap.String("-f", cfg.FileStorage),
		zap.String("-d", cfg.DSN),
	)

	// Load data from file storage if specified
	if cfg.FileStorage != "" {
		err := usecase.Load(cfg.FileStorage)
		if err != nil {
			log.Error("failed to load data from file", zap.Error(err))
			os.Exit(1)
		}
	}

	// Initialize database or in-memory data storage
	if cfg.DSN != "" {
		db, err := usecase.NewDB(cfg.DSN, log)
		if err != nil {
			log.Fatal("failed to connect to database", zap.Error(err))
		}
		defer func() {
			if err := db.Stop(); err != nil {
				log.Error("error stopping DB", zap.Error(err))
			}
		}()
		uc = usecase.New(db, log, cfg)
	} else {
		db := usecase.NewData(log, cfg)
		uc = usecase.New(db, log, cfg)
	}

	// Initialize controllers
	controller := controllers.New(uc, log, cfg)

	// Setup router
	r := chi.NewRouter()
	r.Mount("/", controller.Router(r))

	// Setup and start the server
	srv := setupServer(r)
	log.Info("server starting", zap.String("address", cfg.Host), zap.String("url", cfg.BaseURL))

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", zap.Error(err))
			sigs <- os.Interrupt
		}
	}()
	log.Info("server started")

	// Wait for shutdown signal
	<-sigs
	ctxTime, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelShutdown()

	// Attempt to gracefully shutdown the server
	if err := srv.Shutdown(ctxTime); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}
