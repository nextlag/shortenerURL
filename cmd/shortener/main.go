package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.Cfg
	fmt.Println(cfg)

	router := chi.NewRouter()
	db := storage.NewInMemoryStorage()

	// Добавляем middleware для логирования запросов
	router.Use(middleware.Logger)

	// Создаем маршрут для обработки GET запросов
	router.Get("/{id}", handlers.GetHandler(db))

	// Создаем маршрут для обработки POST запросов
	router.Post("/", handlers.PostHandler(db))

	log.Fatal(http.ListenAndServe(cfg.Address, router))
}
