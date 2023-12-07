package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/nextlag/shortenerURL/internal/service/app"
)

// GetHandler - обработчик GET-запросов для перенаправления на исходный URL. Принимает storage (database) для поиска сокращенных URL.
func GetHandler(db app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметр "id" из URL, который представляет собой сокращенную версию URL
		id := chi.URLParam(r, "id")

		// Пытаемся найти оригинальный URL в хранилище
		originalURL, err := db.Get(r.Context(), id)
		if err != nil {
			http.Error(w, "not found 400", http.StatusBadRequest)
			return
		}

		// Задаем заголовок Location с оригинальным URL
		w.Header().Set("Location", originalURL)

		// Устанавливаем статус HTTP 307 Temporary Redirect и выполняем перенаправление
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
