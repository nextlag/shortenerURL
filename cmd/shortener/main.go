package main

import (
	"github.com/nextlag/shortenerURL/internal/handlers"
	"github.com/nextlag/shortenerURL/internal/storage"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	db := storage.NewInMemoryStorage()      // Создайте экземпляр хранилища
	mux.HandleFunc(`/`, handlers.Route(db)) // Передайте хранилище в обработчик
	log.Fatal(http.ListenAndServe(":8080", mux))
}
