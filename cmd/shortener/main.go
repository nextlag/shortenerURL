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

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/transport/http/middleware/gzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/transport/http/middleware/zaplogger"
	"github.com/nextlag/shortenerURL/internal/transport/http/router"
)

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    app.New().Cfg.Address, // Получение адреса из настроек
		Handler: router,
	}
}

func main() {
	if err := config.InitializeArgs(); err != nil {
		log.Fatal(err)
	}

	logger := app.New().Log // Создание и настройка логгера
	flag.Parse()            // Парсинг флагов командной строки

	// Создание хранилища данных в памяти
	db := app.New().Stor
	err := db.Load(app.New().Cfg.FileStorage)
	if err != nil {
		_ = fmt.Errorf("failed to load data from file: %v", err)
	}

	// Создание и настройка маршрутов и HTTP-сервера
	rout := router.SetupRouter(db, logger)
	// middleware для логирования запросов
	chi.NewRouter().Use(mwLogger.New(logger))
	mw := gzip.New(rout.ServeHTTP)
	srv := setupServer(mw)

	logger.Info("server starting", zap.String("address", app.New().Cfg.Address), zap.String("service", app.New().Cfg.URLShort))

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
