package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
	"github.com/nextlag/shortenerURL/internal/transport/rest/middleware/gzip"
	"github.com/nextlag/shortenerURL/internal/transport/rest/router"
)

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    app.New().Cfg.Host, // Получение адреса из настроек
		Handler: router,
	}
}

func main() {
	if err := config.MakeConfig(); err != nil {
		log.Fatal(err)
	}
	flag.Parse() // Парсинг флагов командной строки

	var run = app.New()
	var log = run.Log
	var cfg = run.Cfg
	var db app.Storage

	if cfg.DSN != "" {
		stor, err := dbstorage.New(cfg.DSN, log)
		if err != nil {
			log.Error("failed to connect in database", zap.Error(err))
		}
		defer func() {
			if err := stor.Stop(); err != nil {
				log.Error("error stopping DB", zap.Error(err))
			}
		}()
		db = stor
	} else {
		db = storage.New(log, cfg)
	}

	log.Info("initialized flags",
		zap.String("-a", cfg.Host),
		zap.String("-b", cfg.BaseURL),
		zap.String("-f", cfg.FileStorage),
		zap.String("-d", cfg.DSN),
	)

	// Создание хранилища данных в памяти
	if cfg.FileStorage != "" {
		data, err := storage.Load(cfg.FileStorage)
		if err != nil {
			log.Error("failed to load data from file", zap.Error(err))
			os.Exit(1) // Завершаем программу при ошибке загрузки данных
		}
		db = data
	}

	// Создание и настройка маршрутов и HTTP-сервера
	rout := router.SetupRouter(db, log, cfg)
	mw := gzip.New(rout.ServeHTTP)
	srv := setupServer(mw)

	log.Info("server starting",
		zap.String("address", cfg.Host),
		zap.String("url", cfg.BaseURL))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартовал, логируем ошибку
			log.Error("failed to start server", zap.String("error", err.Error()))
			done <- os.Interrupt
		}
	}()
	log.Info("server started")

	<-done // Ожидание сигнала завершения
	log.Info("server stopped")
}
