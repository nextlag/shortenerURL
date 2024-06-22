package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nextlag/shortenerURL/internal/config"
)

func TestResponseConflict(t *testing.T) {
	config.Cfg = config.HTTPServer{
		Host: "http://localhost",
	}

	w := httptest.NewRecorder()
	ResponseConflict(w, "abc123")

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	expectedResult := "/abc123"
	if response.Result != expectedResult {
		t.Errorf("expected result %s, got %s", expectedResult, response.Result)
	}
}

func TestResponseCreated(t *testing.T) {
	config.Cfg = config.HTTPServer{
		Host: "http://localhost",
	}

	w := httptest.NewRecorder()
	ResponseCreated(w, "xyz789")

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	expectedResult := "/xyz789"
	if response.Result != expectedResult {
		t.Errorf("expected result %s, got %s", expectedResult, response.Result)
	}
}
