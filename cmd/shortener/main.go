package main

import (
	"fmt"
	"log"
	"net/http"
)

type database map[string]string

var db = database{
	"EwHXdJfB": "https://practicum.yandex.ru/",
}

func route(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		db.getHandler(w, r)
	case http.MethodPost:
		db.postHandler(w, r)
	default:
		fmt.Errorf("bad request 400", http.StatusBadRequest)
	}
}

func (db database) getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Path[1:]
		// Проверяем, есть ли соответствующий оригинальный URL для данного идентификатора
		originalURL, ok := db[id]
		if !ok {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		// Location = оригинальный URL
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "bad request 400", http.StatusBadRequest)
}

func (db database) postHandler(w http.ResponseWriter, r *http.Request) {
	// Заголовки ответа
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "30")
	// 201 status
	w.WriteHeader(http.StatusCreated)
	// Отправляем тело ответа
	_, err := fmt.Fprintf(w, "http://localhost:8080/EwHXdJfB")
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, route)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
