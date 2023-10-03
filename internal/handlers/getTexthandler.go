package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

func GetHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		originalURL, ok := db.Get(id)
		if !ok {
			http.Error(w, "bad request 400", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
