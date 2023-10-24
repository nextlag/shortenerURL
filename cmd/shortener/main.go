package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers/httpServer"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/zaplogger"
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
	// создаем роутер
	router := chi.NewRouter()
	mw := mwLogger.New(log)
	// Настройка обработчиков маршрутов для GET и POST запросов
	router.With(mw).Get("/{id}", httpServer.GetHandler(db))
	router.With(mw).Post("/api/shorten", httpServer.Shorten(log, db))
	router.With(mw).Post("/", httpServer.Save(db))
	return router
}

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Args.Address,
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
	defer logger.Sync()
	return logger
}

func main() {
	logger := setupLogger()
	flag.Parse()

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Настройка маршрутов
	rout := setupRouter(db, logger)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))

	// Создание HTTP-сервера с настроенными маршрутами
	srv := setupServer(rout)

	logger.Info("server starting", zap.String("address", config.Args.Address), zap.String("url", config.Args.URLShort))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартанул вернуть ошибку
			logger.Error("failed to start server", zap.String("error", err.Error()))
			done <- os.Interrupt
		}
	}()
	logger.Info("server started")
	<-done
	logger.Info("server stopped")
}
