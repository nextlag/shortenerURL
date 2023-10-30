package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

// GetHandler - обработчик GET-запросов для перенаправления на исходный Url. Принимает storage (db) для поиска сокращенных Url.
func GetHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметр "alias" из Url, который представляет собой сокращенную версию Url
		alias := chi.URLParam(r, "alias")

		// Пытаемся найти оригинальный Url в хранилище
		url, err := db.Get(alias)
		if err != nil {
			http.Error(w, "not found 400", http.StatusBadRequest)
			return
		}

		// Задаем заголовок Location с оригинальным Url
		w.Header().Set("Location", url)

		// Устанавливаем статус HTTP 307 Temporary Redirect и выполняем перенаправление
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
