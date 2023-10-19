package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers"
	mwLogger "github.com/nextlag/shortenerURL/internal/middleware/logger"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"log/slog"
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

func setupRouter(db storage.Storage, log *slog.Logger) *chi.Mux {
	// создаем роутер
	router := chi.NewRouter()
	mw := mwLogger.New(log)

	// Настройка обработчиков маршрутов для GET и POST запросов
	router.With(mw).Get("/{id}", handlers.GetHandler(db))
	router.With(mw).Post("/", handlers.PostHandler(db))
	return router
}

func setupServer(router http.Handler) *http.Server {
	// Создание HTTP-сервера с указанным адресом и обработчиком маршрутов
	return &http.Server{
		Addr:    config.Args.Address,
		Handler: router,
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func main() {
	log := setupLogger()
	flag.Parse()

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Настройка маршрутов
	rout := setupRouter(db, log)
	router := chi.NewRouter()
	router.Use(mwLogger.New(log))

	// Создание HTTP-сервера с настроенными маршрутами
	srv := setupServer(rout)

	log.Info("server starting", slog.String("address", config.Args.Address), slog.String("url", config.Args.URLShort))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск HTTP-сервера
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			// Если сервер не стартанул вернуть ошибку
			log.Error("failed to start server", slog.String("error", err.Error()))
			done <- os.Interrupt
		}
	}()

	log.Info("server started")
	<-done
	log.Info("server stopped")
}
