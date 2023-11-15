package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/database/dbstorage"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/transport/rest/middleware/gzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/transport/rest/middleware/zaplogger"
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

	logger := app.New().Log // Создание и настройка логгера

	if config.Config.DSN != "" {
		db, err := dbstorage.New(config.Config.DSN)
		if err != nil {
			logger.Error("failed to connect in database", zap.Error(err))
		}
		// Закрытие соединения с базой данных при завершении работы
		defer func() {
			if err := db.Stop(); err != nil {
				logger.Error("error stopping DB", zap.Error(err))
			}
		}()
	}

	// flag.Parse() // Парсинг флагов командной строки

	fmt.Printf("Address: %s\n", config.Config.Address)
	fmt.Printf("URLShort: %s\n", config.Config.URLShort)
	fmt.Printf("FileStorage: %s\n", config.Config.FileStorage)
	fmt.Printf("DSN: %s\n", config.Config.DSN)
	// Создание хранилища данных в памяти
	stor := app.New().Stor
	err := stor.Load(app.New().Cfg.FileStorage)
	if err != nil {
		_ = fmt.Errorf("failed to load data from file: %v", err)
	}

	// Создание и настройка маршрутов и HTTP-сервера
	rout := router.SetupRouter(stor, logger)
	// middleware для логирования запросов
	chi.NewRouter().Use(mwLogger.New(logger))
	mw := gzip.New(rout.ServeHTTP)
	srv := setupServer(mw)

	logger.Info("server starting",
		zap.String("address", app.New().Cfg.Address),
		zap.String("url", app.New().Cfg.URLShort))

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
