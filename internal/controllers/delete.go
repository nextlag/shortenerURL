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
// and initiates the use case layer to perform the deletion asynchronously.
func (c *Controller) Del(w http.ResponseWriter, r *http.Request) {
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		http.Error(w, "You have no links to delete", http.StatusUnauthorized)
		return
	}

	var URLs []string
	if err = json.NewDecoder(r.Body).Decode(&URLs); err != nil {
		c.log.Error("Failed to read json: ", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	go c.deleteURLs(uuid, URLs)

	// Responding to the client immediately
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"aliases sent for deletion": URLs,
	}
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		c.log.Error("Failed to write response: ", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// deleteURLs performs the deletion of URLs
func (c *Controller) deleteURLs(uuid int, URLs []string) {
	for _, url := range URLs {
		c.uc.DoDel(uuid, []string{url}) // Passing nil as context
	}
}
