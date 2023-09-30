package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	db := storage.NewInMemoryStorage()

	// Добавляем middleware для логирования запросов
	r.Use(middleware.Logger)

	// Создаем маршрут для обработки GET запросов
	r.Get("/{id}", handlers.GetHandler(db))

	// Создаем маршрут для обработки POST запросов
	r.Post("/", handlers.PostHandler(db))

	log.Fatal(http.ListenAndServe(":8080", r))
}
