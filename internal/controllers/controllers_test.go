package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/controllers/mocks"
	"github.com/nextlag/shortenerURL/internal/middleware/logger"
	"github.com/nextlag/shortenerURL/internal/usecase"
)

func Ctrl(t *testing.T) (*Controller, *mocks.MockUseCase, *usecase.UseCase) {
	t.Helper()
	l := logger.SetupLogger()
	var cfg config.HTTPServer
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	db := mocks.NewMockUseCase(mockCtl)
	repo := usecase.NewMockRepository(mockCtl)
	uc := usecase.New(repo, l, cfg)
	controller := New(db, l, cfg)
	return controller, db, uc
}

func TestController(t *testing.T) {
	ctrl, db, _ := Ctrl(t)
	r := chi.NewRouter()
	ctrl.Controller(r)

	tests := []struct {
		name           string
		method         string
		target         string
		body           string
		contentType    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:           "GET /{id}",
			method:         http.MethodGet,
			target:         "/testid",
			expectedStatus: http.StatusTemporaryRedirect,
			mockSetup: func() {
				db.EXPECT().DoGet(gomock.Any(), "testid").Return("http://example.com", false, nil).Times(1)
			},
		},
		{
			name:           "GET /ping",
			method:         http.MethodGet,
			target:         "/ping",
			expectedStatus: http.StatusOK,
			mockSetup: func() {
				db.EXPECT().DoHealthcheck().Return(true, nil).Times(1)
			},
		},
		{
			name:           "POST /api/shorten",
			method:         http.MethodPost,
			target:         "/api/shorten",
			body:           `{"url": "http://example.com"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusCreated,
			mockSetup: func() {
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any()).Return("shortened_url", nil).Times(1)
			},
		},
		{
			name:           "POST /",
			method:         http.MethodPost,
			target:         "/",
			body:           "http://example.com",
			contentType:    "text/plain",
			expectedStatus: http.StatusCreated,
			mockSetup: func() {
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any()).Return("shortened_url", nil).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
