package handlers

import (
	"net/http"

	"github.com/nextlag/shortenerURL/internal/storage/database/dbstorage"
)

type PingHandler struct {
	db *dbstorage.DBStorage
}

func NewHealCheck(db *dbstorage.DBStorage) *PingHandler {
	return &PingHandler{
		db: db,
	}
}

func (h *PingHandler) healCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.db == nil || !h.db.CheckConnection() {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
