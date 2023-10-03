package handlers

import (
	"github.com/nextlag/shortenerURL/internal/storage"
	"net/http"
)

func Route(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "text/plain":
			TEXTPostHandler(db)(w, r)
		case "application/json":
			JSONPostHandler(db)(w, r)
		default:
			http.Error(w, "unsupported media type 400", http.StatusBadRequest)
		}
	}
}
