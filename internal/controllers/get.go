// Package controllers provides HTTP handlers for URL shortening service.
package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Get handles GET requests for redirecting to the original URL. Accepts a storage (database) to look up shortened URLs.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" parameter from the URL, which represents the shortened version of the URL.
	id := chi.URLParam(r, "id")

	// Attempt to find the original URL in the storage.
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

	// Set the Location header to the original URL.
	w.Header().Set("Location", originalURL)

	// Set the HTTP status to 307 Temporary Redirect and perform the redirect.
	w.WriteHeader(http.StatusTemporaryRedirect)
}
