package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
)

type DelURL struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

func NewDelURL(db app.Storage, log *zap.Logger, cfg config.Args) *DelURL {
	return &DelURL{db: db, log: log, cfg: cfg}
}

func (h *DelURL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid, err := auth.CheckCookie(w, r, h.log)
	if err != nil {
		h.log.Error("Error getting cookie: ", zap.Error(err))
		w.WriteHeader(401)
		w.Write([]byte("You have no links to delete"))
		return
	}

	var URLs []string

	err = json.NewDecoder(r.Body).Decode(&URLs)
	if err != nil {
		h.log.Error("Failed to read json: ", zap.Error(err))
	}

	h.db.Del(r.Context(), uuid, URLs)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(202)
	w.Write([]byte("Deleted"))

}
