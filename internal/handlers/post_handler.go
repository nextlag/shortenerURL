package handlers

import (
	"fmt"
	"github.com/nextlag/shortenerURL/internal/storage"
	"io"
	"math/rand"
	"net/http"
)

func PostHandler(db *storage.InMemoryStorage, w http.ResponseWriter, r *http.Request) {
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
	db.Put(shortURL, string(body))
}

func generateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
