// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Get handles GET requests for redirecting to the original URL. It takes a storage (database) to lookup shortened URLs.
// It extracts the "id" parameter from the URL, searches for the original URL in the storage, and redirects to it.
// If the URL is marked as deleted, it returns a 410 Gone status.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" parameter from the URL, which represents the shortened version of the URL
	id := chi.URLParam(r, "id")

	// Try to find the original URL in the storage
	originalURL, deleteStatus, err := c.uc.DoGet(r.Context(), id)
	if err != nil {
		http.Error(w, "not found 400", http.StatusBadRequest)
		return
	}
	if deleteStatus {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Deleted URL"))
		return
	}

	// Set the Location header to the original URL
	w.Header().Set("Location", originalURL)

	// Set HTTP status 307 Temporary Redirect and perform the redirection
	w.WriteHeader(http.StatusTemporaryRedirect)
}
