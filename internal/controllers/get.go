package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Get обработчик GET-запросов для перенаправления на исходный URL. Принимает storage (database) для поиска сокращенных URL.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр "id" из URL, который представляет собой сокращенную версию URL
	id := chi.URLParam(r, "id")

	// Пытаемся найти оригинальный URL в хранилище
	originalURL, deleteStatus, err := c.uc.DoGet(r.Context(), id)
	if err != nil {
		http.Error(w, "not found 400", http.StatusBadRequest)
		return
	}
	if deleteStatus {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Deleted URL"))
		return
	}

	// Задаем заголовок Location с оригинальным URL
	w.Header().Set("Location", originalURL)

	// Устанавливаем статус HTTP 307 Temporary Redirect и выполняем перенаправление
	w.WriteHeader(http.StatusTemporaryRedirect)
}
