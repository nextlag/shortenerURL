package handlers

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
	"github.com/nextlag/shortenerURL/internal/utils/lg"
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
func (h *GetAllURLsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log := lg.New()
	var cfg config.Args

	userID := auth.CheckCookieForID(res, req, log)

	switch userID {
	case -1:
		res.WriteHeader(401)
		res.Write([]byte("Unauthorized"))
	default:
		userURLs, err := h.db.GetAll(req.Context(), userID, cfg.URLShort)
		if err != nil {
			h.log.Error("Error getting URLs by ID", zap.Error(err))
		}

		res.Header().Set("Content-Type", "application/json")
		if string(userURLs) == "null" {
			res.WriteHeader(204)
			res.Write([]byte("No content"))
		} else {
			res.WriteHeader(201)
			res.Write(userURLs)
		}
	}
}

// func (h *GetAllURLsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	userID := auth.CheckCookieForID(w, r, h.log)
// 	var cfg config.Args
//
// 	switch userID {
// 	case -1:
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	default:
// 		userURLs, err := h.db.GetAll(r.Context(), userID, cfg.URLShort)
// 		if err != nil {
// 			h.log.Error("Error getting URLs by ID", zap.Error(err))
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		if string(userURLs) == "null" {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
//
// 		w.WriteHeader(http.StatusCreated)
// 		w.Header().Set("Content-Type", "application/json")
// 		if string(userURLs) == "null" {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}
//
// 		w.WriteHeader(http.StatusCreated)
// 		err = json.NewEncoder(w).Encode(userURLs)
// 		if err != nil {
// 			h.log.Error("Failed to encode JSON response", zap.Error(err))
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}
// 	}
// }
