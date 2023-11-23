package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/storage/database/dbstorage"
	"github.com/nextlag/shortenerURL/internal/transport/rest/middleware/gzip"
	"github.com/nextlag/shortenerURL/internal/transport/rest/router"
)

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    app.New().Cfg.Address, // Получение адреса из настроек
		Handler: router,
	}
}

func main() {
	if err := config.MakeConfig(); err != nil {
		log.Fatal(err)
	}
	flag.Parse() // Парсинг флагов командной строки

	var logger = app.New().Log // Создание и настройка логгера
	var cfg = app.New().Cfg
	var stor app.Storage

	if cfg.DSN != "" {
		db, err := dbstorage.New(config.Config.DSN)
		if err != nil {
			logger.Error("failed to connect in database", zap.Error(err))
		}
		defer func() {
			if err := db.Stop(); err != nil {
				logger.Error("error stopping DB", zap.Error(err))
			}
		}()
		stor = db
	} else {
		stor = storage.New()
	}

	logger.Info("initialized flags",
		zap.String("-a", cfg.Address),
		zap.String("-b", cfg.URLShort),
		zap.String("-f", cfg.FileStorage),
		zap.String("-d", cfg.DSN),
	)

	// Создание хранилища данных в памяти
	// stor = app.New().Stor
	if cfg.FileStorage != "" {
		// err := storage.Load(cfg.FileStorage)
		err := storage.New().Load(cfg.FileStorage)
		if err != nil {
			_ = fmt.Errorf("failed to load data from file: %v", err)
		}
	}

	// Создание и настройка маршрутов и HTTP-сервера
	rout := router.SetupRouter(stor, logger)
	mw := gzip.New(rout.ServeHTTP)
	srv := setupServer(mw)

	logger.Info("server starting",
		zap.String("address", cfg.Address),
		zap.String("url", cfg.URLShort))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартовал, логируем ошибку
			logger.Error("failed to start server", zap.String("error", err.Error()))
			done <- os.Interrupt
		}
	}()
	logger.Info("server started")

	<-done // Ожидание сигнала завершения
	logger.Info("server stopped")
}
