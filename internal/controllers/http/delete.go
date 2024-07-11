// Package controllers provides the handlers for managing URL shortening operations.
package http

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

	var aliases []string
	if err = json.NewDecoder(r.Body).Decode(&aliases); err != nil {
		c.log.Error("Failed to read json: ", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Запуск удаления URL в горутине
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for _, alias := range aliases {
			c.uc.DoDel(uuid, []string{alias})
		}
	}()

	// ответ клиенту, какие алиасы отправлены на удаление
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"aliases sent for deletion": aliases,
	}
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		c.log.Error("Failed to write response: ", zap.Error(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
