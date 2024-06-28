// Package controllers provides the handlers for managing URL shortening operations.
package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

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
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("You have no links to delete"))
		return
	}

	var URLs []string

	err = json.NewDecoder(r.Body).Decode(&URLs)
	if err != nil {
		c.log.Error("Failed to read json: ", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}

	type result struct {
		URL string
		Err error
	}

	results := make(chan result, len(URLs))

	var wg sync.WaitGroup
	for _, url := range URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err = c.uc.DoDel(r.Context(), uuid, []string{url})
			results <- result{URL: url, Err: err}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var (
		successfulURLs []string
		failedURLs     []string
	)
	for res := range results {
		if res.Err != nil {
			c.log.Error("Error deleting user URL", zap.String("url", res.URL), zap.Error(res.Err))
			failedURLs = append(failedURLs, res.URL)
		} else {
			successfulURLs = append(successfulURLs, res.URL)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"deleted": successfulURLs,
	}

	if len(failedURLs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = "Some URLs could not be deleted"
		response["failed"] = failedURLs
	} else {
		w.WriteHeader(http.StatusAccepted)
		response["message"] = "All URLs deleted successfully"
	}

	json.NewEncoder(w).Encode(response)
}
