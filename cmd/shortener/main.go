package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	router := chi.NewRouter()
	db := storage.NewInMemoryStorage()

	// Создаем маршрут для обработки GET запросов
	router.Get("/{id}", handlers.GetHandler(db))

	// Создаем маршрут для обработки POST запросов
	router.Post("/", handlers.PostHandler(db))

	log.Printf("Start server: %s | ShortenerURL: %s", *config.Address, *config.URLShort)
	log.Fatal(http.ListenAndServe(*config.Address, router))
}
