package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		w.WriteHeader(401)
		w.Write([]byte("You have no links to delete"))
		return
	}

	var URLs []string

	err = json.NewDecoder(r.Body).Decode(&URLs)
	if err != nil {
		c.log.Error("Failed to read json: ", zap.Error(err))
	}

	err = c.uc.DoDel(r.Context(), uuid, URLs)
	if err != nil {
		c.log.Error("error deleting user url", zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(202)
	w.Write([]byte("Deleted"))

}
