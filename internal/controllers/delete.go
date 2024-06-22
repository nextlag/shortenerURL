// Package controllers provides HTTP handlers for URL shortening service.
package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// Del handles an HTTP request to delete multiple URLs associated with a user.
func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	// Check the authentication cookie.
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("You have no links to delete"))
		return
	}

	var URLs []string

	// Decode the JSON request body.
	err = json.NewDecoder(r.Body).Decode(&URLs)
	if err != nil {
		c.log.Error("Failed to read json: ", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to decode request body"))
		return
	}

	// Delete the URLs.
	err = c.uc.DoDel(r.Context(), uuid, URLs)
	if err != nil {
		c.log.Error("error deleting user URL", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete URLs"))
		return
	}

	// Set the response headers and write the response.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Deleted"))
}
