package handlers

import (
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

func GetHandler(db *storage.InMemoryStorage, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Получаем идентификатор
		id := r.URL.Path[1:]
		// Проверяем, есть ли соответствующий оригинальный URL для данного идентификатора
		originalURL, ok := db.Get(id)
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
