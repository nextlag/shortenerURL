package controllers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/controllers"
)

type mockUsecase struct{}

func (m *mockUsecase) DoGet(ctx context.Context, alias string) (string, bool, error) {
	return "http://example.com", false, nil
}

func (m *mockUsecase) DoGetAll(ctx context.Context, userID int, url string) ([]byte, error) {
	return []byte(`[{"short_url": "http://example.com/short1", "original_url": "http://example.com/original1"}]`), nil
}

func (m *mockUsecase) DoPut(ctx context.Context, url string, uuid int) (string, error) {
	return "shortened_url", nil
}

func (m *mockUsecase) DoDel(ctx context.Context, id int, aliases []string) error {
	return nil
}

func (m *mockUsecase) DoHealthcheck() (bool, error) {
	return true, nil
}

// ExampleController_HealthCheck demonstrates how to use the HealthCheck endpoint.
func ExampleController_HealthCheck() {
	// Setup
	log := zap.NewNop()
	cfg := config.HTTPServer{}
	uc := &mockUsecase{}

	ctrl := controllers.New(uc, log, cfg)
	r := chi.NewRouter()
	ctrl.Router(r)

	// Example request to the /ping endpoint
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Output response status
	if res.StatusCode == http.StatusOK {
		println("Healthcheck OK")
	} else {
		println("Healthcheck Failed")
	}

	// Output:
}

// TestController_HealthCheck tests the HealthCheck endpoint of the Controller.
func TestController_HealthCheck(t *testing.T) {
	log := zap.NewNop()
	cfg := config.HTTPServer{}
	uc := &mockUsecase{}

	ctrl := controllers.New(uc, log, cfg)
	r := chi.NewRouter()

	// Handle not found routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})

	// Handle method not allowed
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "400 bad request", http.StatusBadRequest)
	})

	// Setup routes with the controller
	ctrl.Router(r)

	tests := []struct {
		method     string
		target     string
		statusCode int
	}{
		{method: http.MethodGet, target: "/ping", statusCode: http.StatusOK},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.target, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := w.Result()
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}(res.Body)

		if res.StatusCode != test.statusCode {
			t.Errorf("expected status %d, got %d", test.statusCode, res.StatusCode)
		}
	}
}
