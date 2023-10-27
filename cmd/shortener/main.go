package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers/httpserver"
	"github.com/nextlag/shortenerURL/internal/middleware/mwGzip"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/mwZapLogger"
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

func setupRouter(db storage.Storage, log *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // добавляем уникальный идентификатор
	router.Use(middleware.Logger)    // добавляем вывод стандартного логгера

	// Создание экземпляра middleware.Logger
	mw := mwLogger.New(log)

	// Настройка маршрутов с использованием middleware
	router.With(mw).Get("/{id}", httpserver.GetHandler(db))
	router.With(mw).Post("/api/shorten", httpserver.Shorten(log, db))
	router.With(mw).Get("/api/shorten/{id}", httpserver.GetHandler(db))
	router.With(mw).Post("/", httpserver.Save(db))

	return router
}

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Args.Address, // Получение адреса из настроек
		Handler: router,
	}
}

func setupLogger() *zap.Logger {
	// Настраиваем конфигурацию логгера
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // Уровень логирования

	// Создаем логгер
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Отложенное закрытие логгера
	return logger
}

func main() {
	logger := setupLogger() // Создание и настройка логгера
	flag.Parse()            // Парсинг флагов командной строки

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Создание и настройка маршрутов и HTTP-сервера
	router := setupRouter(db, logger)
	// middleware для логирования запросов
	chi.NewRouter().Use(mwLogger.New(logger))
	mw := mwGzip.New(router.ServeHTTP)
	srv := setupServer(mw)
	//srv := setupServer(router)

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
