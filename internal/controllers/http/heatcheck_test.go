package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/configuration"
	http2 "github.com/nextlag/shortenerURL/internal/controllers/http"
	"github.com/nextlag/shortenerURL/internal/entity"
)

type mockUsecase struct{}

func (m *mockUsecase) DoGet(ctx context.Context, alias string) (*entity.URL, error) {
	return &entity.URL{URL: "http://example.com", Alias: alias, IsDeleted: false}, nil
}

func (m *mockUsecase) DoGetAll(ctx context.Context, userID int, url string) ([]*entity.URL, error) {
	return []*entity.URL{
		{Alias: "short1", URL: "http://example.com/original1"},
	}, nil
}

func (m *mockUsecase) DoPut(ctx context.Context, url string, alias string, uuid int) (string, error) {
	return "shortened_url", nil
}

func (m *mockUsecase) DoDel(ctx context.Context, id int, aliases []string) {}

func (m *mockUsecase) DoHealthcheck() (bool, error) {
	return true, nil
}

func (m *mockUsecase) DoGetStats(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func ExampleController_HealthCheck() {
	log := zap.NewNop()
	cfg := configuration.Config{}
	uc := &mockUsecase{}

	wg := sync.WaitGroup{}
	ctrl := http2.New(uc, &wg, &cfg, log)
	r := chi.NewRouter()
	ctrl.Controller(r)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		println("Healthcheck OK")
	} else {
		println("Healthcheck Failed")
	}
}

func TestController_HealthCheck(t *testing.T) {
	log := zap.NewNop()
	cfg := configuration.Config{}
	uc := &mockUsecase{}
	wg := sync.WaitGroup{}

	ctrl := http2.New(uc, &wg, &cfg, log)
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "400 bad request", http.StatusBadRequest)
	})

	ctrl.Controller(r)

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
		defer res.Body.Close()
		if res.StatusCode != test.statusCode {
			t.Errorf("expected status %d, got %d", test.statusCode, res.StatusCode)
		}
	}
}
