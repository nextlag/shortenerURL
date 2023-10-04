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

func init() {
	if err := config.InitializeArgs(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	router := chi.NewRouter()
	db := storage.NewInMemoryStorage()

	// Создаем маршрут для обработки GET запросов
	router.Get("/{id}", handlers.GetHandler(db))

	// Создаем маршрут для обработки POST запросов
	router.Post("/", handlers.PostHandler(db))

	log.Printf("START HTTPServer: %s | ShortenerURL: %s", config.AFlag.Address, config.AFlag.URLShort)
	log.Fatal(http.ListenAndServe(config.AFlag.Address, router))
}
