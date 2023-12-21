package handlers

import (
	"net/http"

	"github.com/nextlag/shortenerURL/internal/service/app"
)

type HealtCheck struct {
	db app.Storage
}

func NewHealtCheck(db app.Storage) *HealtCheck {
	return &HealtCheck{
		db: db,
	}
}

func (h *HealtCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.db == nil || !h.db.Healtcheck() {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
