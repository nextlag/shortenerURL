package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
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
	var cfg config.Args

	userID := auth.CheckCookie(w, r, h.log)

	switch userID {
	case -1:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	default:
		allURL, err := h.db.GetAll(r.Context(), userID, cfg.URLShort)
		if err != nil {
			h.log.Error("Error getting URLs by ID", zap.Error(err))
		}

		w.Header().Set("Content-Type", "application/json")
		if string(allURL) == "null" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("No content"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(allURL)
		}
	}
}
