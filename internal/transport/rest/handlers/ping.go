package handlers

import (
	"net/http"

	"github.com/nextlag/shortenerURL/internal/service/app"
)

type PingHandler struct {
	db app.Storage
}

func NewHealCheck(db app.Storage) *PingHandler {
	return &PingHandler{
		db: db,
	}
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.db == nil || !h.db.CheckConnection(r.Context()) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
