// Package controllers provides the handlers for managing URL shortening operations.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Get handles GET requests for redirecting to the original URL.
// It extracts the "id" parameter from the URL, searches for the original URL in the storage,
// and redirects to it. If the URL is marked as deleted, it returns a 410 Gone status.
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := c.uc.DoGet(r.Context(), id)
	if err != nil {
		c.log.Error("error", zap.Error(err))
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if url.IsDeleted {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Deleted URL"))
		return
	}

	w.Header().Set("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
