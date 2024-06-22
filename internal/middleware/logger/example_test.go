package logger_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
)

// Example of using the logging middleware in an HTTP server.
func ExampleNew() {
	// Initialize the logger
	log := logger.SetupLogger()

	// Define a simple handler that returns a plain text response.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, logged world!"))
	})

	// Create a new router and add the logging middleware.
	r := chi.NewRouter()
	r.Use(logger.New(log, config.HTTPServer{
		FileStorage: "/tmp/data",
	}))

	// Add the handler to the router.
	r.Get("/", handler)

	// Create a new HTTP request.
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// Serve the HTTP request using the router with logging middleware.
	r.ServeHTTP(w, req)

	// Get the recorded response.
	res := w.Result()
	defer res.Body.Close()

	// Print the response status and body.
	fmt.Println("Status:", res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	fmt.Println("Body:", string(body))

	// Output:
	// Status: 200
	// Body: Hello, logged world!
}
