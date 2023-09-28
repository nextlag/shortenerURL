package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// БД для сокращенных URL
type database map[string]string

var db = database{}

func route(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		db.getHandler(w, r)
	case http.MethodPost:
		db.postHandler(w, r)
	default:
		http.Error(w, "bad request 400", http.StatusBadRequest)
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
	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}

	// Генерируем случайную строку для сокращенного URL
	shortURL := generateRandomString(8)

	// Заголовки ответа
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "30")
	// 201 status
	w.WriteHeader(http.StatusCreated)
	// Отправляем тело ответа с сокращенным URL
	_, err = fmt.Fprintf(w, "http://localhost:8080/%s", shortURL)
	if err != nil {
		return
	}

	// Сохраняем сокращенный URL в базу данных
	db[shortURL] = string(body)
}

func generateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, route)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
