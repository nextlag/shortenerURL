package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
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

func setupRouter(db storage.Storage) *chi.Mux {
	// создаем роутер
	router := chi.NewRouter()
	//mw := mwLogger.New(log)
	// Настройка обработчиков маршрутов для GET и POST запросов
	//router.With(mw).Get("/{id}", handlers.GetHandler(db))
	//router.With(mw).Post("/", handlers.PostHandler(db))
	router.Get("/{id}", handlers.GetHandler(db))
	router.Post("/", handlers.PostHandler(db))
	return router
}

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Args.Address,
		Handler: router,
	}
}

//func setupLogger() *zap.Logger {
//	// Настраиваем конфигурацию логгера
//	cfg := zap.NewDevelopmentConfig()
//	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel) // Уровень логирования
//
//	// Создаем логгер
//	logger, err := cfg.Build()
//	if err != nil {
//		panic(err)
//	}
//	return logger
//}

func main() {
	//logger := setupLogger()
	flag.Parse()

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Настройка маршрутов
	rout := setupRouter(db)
	//router := chi.NewRouter()
	//router.Use(mwLogger.New(logger))

	// Создание HTTP-сервера с настроенными маршрутами
	srv := setupServer(rout)

	//logger.Info("server starting", zap.String("address", config.Args.Address), zap.String("url", config.Args.URLShort))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартанул вернуть ошибку
			fmt.Println("failed to start server")
			done <- os.Interrupt
		}
	}()
	<-done
}
