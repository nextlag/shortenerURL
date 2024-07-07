package main

import (
	"context"
	"errors"
	"flag"
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
	if _, err := configuration.Load(); err != nil {
		log.Fatal("failed to init configuration", zap.Error(err))
	}

	flag.Parse()
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal("Failed to get configuration")
		return
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
	)

	if cfg.FileStorage != "" {
		err := usecase.Load(cfg.FileStorage)
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", zap.Error(err))
			sigs <- os.Interrupt
		}
	}()
	log.Info("server started")

	<-sigs
	ctxTime, cancelShutdown := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxTime); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}
