package handlers

import (
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

func Route(db *storage.InMemoryStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetHandler(db, w, r)
		case http.MethodPost:
			PostHandler(db, w, r)
		default:
			http.Error(w, "bad request 400", http.StatusBadRequest)
		}
	}
}
