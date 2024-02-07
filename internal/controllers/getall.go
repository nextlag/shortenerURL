package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Unauthorized access : ", zap.Error(err))
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}
	userURLs, err := c.uc.DoGetAll(r.Context(), userID, c.cfg.BaseURL)
	if err != nil {
		c.log.Error("Error getting URLs by ID", zap.Error(err))
	}

	w.Header().Set("Content-Type", "application/json")
	if string(userURLs) == "null" {
		w.WriteHeader(200)
		w.Write([]byte("No content"))
	} else {
		w.WriteHeader(200)
		w.Write(userURLs)
	}
}
