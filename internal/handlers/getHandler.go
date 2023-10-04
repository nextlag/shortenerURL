package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

// GetHandler - обработчик GET-запросов для перенаправления на исходный URL. Принимает storage (db) для поиска сокращенных URL.
func GetHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметр "id" из URL, который представляет собой сокращенную версию URL
		id := chi.URLParam(r, "id")

		// Пытаемся найти оригинальный URL в хранилище
		originalURL, ok := db.Get(id)
		if !ok {
			// Если сокращенный URL не найден, отправляем ошибку 400 Bad Request
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}

		// Задаем заголовок Location с оригинальным URL
		w.Header().Set("Location", originalURL)

		// Устанавливаем статус HTTP 307 Temporary Redirect и выполняем перенаправление
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
