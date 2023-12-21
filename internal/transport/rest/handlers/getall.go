package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
)

type GetAllURLsHandler struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

func NewGetAllHandler(db app.Storage, log *zap.Logger, cfg config.Args) *GetAllURLsHandler {
	return &GetAllURLsHandler{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

func (h *GetAllURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cfg = app.New().Cfg
	uuid, err := auth.CheckCookie(w, r, h.log)
	if err != nil {
		h.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	switch uuid {
	case -1:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	default:
		allURL, err := h.db.GetAll(r.Context(), uuid, cfg.BaseURL)
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
