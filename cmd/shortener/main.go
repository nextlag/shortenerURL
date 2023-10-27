package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/config"
	lg "github.com/nextlag/shortenerURL/internal/logger"
	mwGzip "github.com/nextlag/shortenerURL/internal/middleware/gzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/zaplogger"
	"github.com/nextlag/shortenerURL/internal/rout"
	"github.com/nextlag/shortenerURL/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := config.InitializeArgs(); err != nil {
		log.Fatal(err)
	}
}

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Args.Address, // Получение адреса из настроек
		Handler: router,
	}
}

func main() {
	logger := lg.SetupLogger() // Создание и настройка логгера
	flag.Parse()               // Парсинг флагов командной строки

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Создание и настройка маршрутов и HTTP-сервера
	router := rout.SetupRouter(db, logger)
	// middleware для логирования запросов
	chi.NewRouter().Use(mwLogger.NewLogger(logger))
	mw := mwGzip.NewGzip(router.ServeHTTP)
	srv := setupServer(mw)
	//srv := setupServer(rout)

	logger.Info("server starting", zap.String("address", config.Args.Address), zap.String("url", config.Args.URLShort))

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
