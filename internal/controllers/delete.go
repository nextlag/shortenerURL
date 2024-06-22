// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// Del handles the HTTP request for deleting URLs associated with a user.
// It checks the user's authentication, decodes the request body to get the list of URLs to delete,
// and calls the use case layer to perform the deletion. If any error occurs, it logs the error and responds
// with an appropriate message.
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
