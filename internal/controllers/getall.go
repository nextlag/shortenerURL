// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// GetAll handles the HTTP request for retrieving all URLs associated with a user.
// It checks the user's authentication, retrieves the URLs from the use case layer,
// and responds with the list of URLs in JSON format. If any error occurs, it logs the error and responds
// with an appropriate message.
func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Unauthorized access: ", zap.Error(err))
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
