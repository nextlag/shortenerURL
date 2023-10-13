package main

import (
	"context"
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func init() {
	if err := config.InitializeArgs(); err != nil {
		log.Fatal(err)
	}
}

func setupRouter(db storage.Storage) *chi.Mux {
	// Создание роутера
	router := chi.NewRouter()

	// Настройка обработчиков маршрутов для GET и POST запросов
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

func handleShutdown(srv *http.Server, idleConnsClosed chan struct{}) {
	// Создание канала для ожидания сигнала завершения операции
	sigint := make(chan os.Signal, 1)

	// Регистрация обработчика сигнала завершения (Ctrl+C)
	signal.Notify(sigint, os.Interrupt)

	// Ожидание сигнала завершения
	<-sigint

	// Завершение работы HTTP-сервера с использованием контекста
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	// Закрытие канала, чтобы разблокировать ожидание завершения программы
	close(idleConnsClosed)
}
func main() {
	flag.Parse()

	// Создание хранилища данных в памяти
	db := storage.NewInMemoryStorage()

	// Настройка маршрутов
	router := setupRouter(db)

	// Создание HTTP-сервера с настроенными маршрутами
	srv := setupServer(router)

	// Создание канала для ожидания сигнала завершения
	idleConnsClosed := make(chan struct{})

	// Запуск обработчика завершения в отдельной горутине
	go handleShutdown(srv, idleConnsClosed)

	// Вывод сообщения о старте сервера
	log.Printf("Server address: %s || Base URL: %s", config.Args.Address, config.Args.URLShort)

	// Запуск HTTP-сервера и сравнение err с http.ErrServerClosed
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Если сервер вернул ошибку, вывести сообщение об ошибке и завершить программу
		log.Fatal(http.ListenAndServe(config.Args.Address, router))
	}
	// Ожидание завершения всех соединений перед завершением программы
	<-idleConnsClosed
}
