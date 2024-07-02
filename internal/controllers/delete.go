package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// DeletionRequest represents the data needed to perform URL deletion.
type DeletionRequest struct {
	Req  *http.Request
	UUID int
	URLs []string
}

// DeletionResult holds the result of a deletion attempt.
type DeletionResult struct {
	SuccessfulURLs []string
	FailedURLs     []string
}

// Del handles the HTTP request for deleting URLs associated with a user.
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

	delReq := DeletionRequest{
		Req:  r,
		UUID: uuid,
		URLs: URLs,
	}

	delRes := c.generateResults(delReq)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"aliases sent for deletion": delRes.SuccessfulURLs,
	}

	if len(delRes.FailedURLs) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		response["message"] = "Some URLs could not be deleted"
		response["failed"] = delRes.FailedURLs
	} else {
		w.WriteHeader(http.StatusAccepted)
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		c.log.Error("Failed to write response: ", zap.Error(err))
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// generateResults processes the deletion of URLs and returns the results.
func (c *Controller) generateResults(delReq DeletionRequest) DeletionResult {
	type result struct {
		URL string
		Err error
	}

	results := make(chan result, len(delReq.URLs))
	var wg sync.WaitGroup

	for _, url := range delReq.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err := c.uc.DoDel(delReq.Req.Context(), delReq.UUID, []string{url})
			results <- result{URL: url, Err: err}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var successfulURLs, failedURLs []string
	for res := range results {
		if res.Err != nil {
			c.log.Error("Error deleting user URL", zap.String("url", res.URL), zap.Error(res.Err))
			failedURLs = append(failedURLs, res.URL)
		} else {
			successfulURLs = append(successfulURLs, res.URL)
		}
	}

	return DeletionResult{
		SuccessfulURLs: successfulURLs,
		FailedURLs:     failedURLs,
	}
}
