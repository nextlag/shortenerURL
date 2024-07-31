package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/cert"
	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/internal/controllers"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	log := logger.SetupLogger()

	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal("failed to init configuration", zap.Error(err))
	}

	fmt.Printf(
		"Build version: %s,\nBuild date: %s,\nBuild commit: %s,\n",
		buildVersion,
		buildDate,
		buildCommit,
	)
	log.Debug(
		"initialized flags",
		zap.String("-a", cfg.Host),
		zap.String("-b", cfg.BaseURL),
		zap.String("-f", cfg.FileStorage),
		zap.String("-d", cfg.DSN),
		zap.String("-c", cfg.ConfigPath),
		zap.Bool("-s", cfg.EnableHTTPS),
		zap.String("-t", cfg.TrustedSubnet),
	)

	if cfg.FileStorage != "" {
		err = usecase.Load(cfg.FileStorage)
		if err != nil {
			log.Fatal("failed to load data from file", zap.Error(err))
		}
	}

	var uc *usecase.UseCase
	if cfg.DSN != "" {
		db, err := usecase.NewDB(cfg.DSN, log)
		if err != nil {
			log.Fatal("failed to connect in database", zap.Error(err))
		}
		defer func() {
			if err = db.Stop(); err != nil {
				log.Error("error stopping DB", zap.Error(err))
			}
		}()
		uc = usecase.New(db, log, cfg)
	} else {
		db := usecase.NewData(log, cfg)
		uc = usecase.New(db, log, cfg)
	}

	controller := controllers.New(uc, log, cfg)

	r := chi.NewRouter()
	r.Mount("/", controller.Controller(r))

	srv := &http.Server{
		Addr:    cfg.Host,
		Handler: r,
	}

	if cfg.EnableHTTPS {
		srv = &http.Server{
			Addr:      cfg.Host,
			Handler:   r,
			TLSConfig: cert.NewCert("localhost").TLSConfig(),
		}
	}

	log.Info(
		"server starting",
		zap.String("address", cfg.Host),
		zap.String("url", cfg.BaseURL),
	)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		log.Info("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err = srv.Shutdown(ctx); err != nil {
			log.Error("HTTP server Shutdown:", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	switch {
	case cfg.EnableHTTPS:
		if err = srv.ListenAndServeTLS(cfg.Cert, cfg.Key); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("HTTPS server ListenAndServeTLS:", zap.Error(err))
		}
	default:
		if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("HTTP server ListenAndServe:", zap.Error(err))
		}
	}

	<-idleConnsClosed

	log.Info("Server Shutdown gracefully")
}
