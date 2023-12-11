package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/service/auth"

	"github.com/nextlag/shortenerURL/internal/service/app"
)

type GetAllURLsHandler struct {
	db  app.Storage
	log *zap.Logger
	ctx context.Context
}

func NewGetAllHandler(log *zap.Logger, db app.Storage) *GetAllURLsHandler {
	return &GetAllURLsHandler{
		log: log,
		db:  db,
	}
}

func (h *GetAllURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cfg := app.New().Cfg
	userID := auth.CheckCookie(w, r, h.log)

	switch userID {
	case -1:
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	default:
		allURL, err := h.db.GetAll(r.Context(), userID, cfg.URLShort)
		if err != nil {
			h.log.Error("Error getting URLs by ID", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if string(allURL) == "null" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No content"))
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(allURL)
		}
	}
}
