package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

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

func TestGetHandler(t *testing.T) {
	tests := []struct {
		Name             string
		RequestPath      string
		ExpectedStatus   int
		ExpectedLocation string
	}{
		{
			Name:             "Valid ID",
			RequestPath:      "/example",
			ExpectedStatus:   http.StatusTemporaryRedirect,
			ExpectedLocation: "http://example.com",
		},
		{
			Name:             "Invalid ID",
			RequestPath:      "/nonexistent",
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedLocation: "",
		}}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl, db, _ := Ctrl(t)

			if test.Name == "Valid ID" {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return("", false, nil).Times(1)
			} else {
				db.EXPECT().DoGet(gomock.Any(), gomock.Any()).Return("", false, errors.New("error")).Times(1)
			}
			// Создаем фейковый запрос
			r := httptest.NewRequest("GET", test.RequestPath, nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(ctrl.Get)
			handler(w, r)

			// Создаем и вызываем handler для маршрута
			resp := w.Result()

			// Проверяем статус кода
			assert.Equal(t, test.ExpectedStatus, resp.StatusCode)
			// Получаем значение Location
			location := resp.Header.Get("Location")
			assert.Empty(t, location)

			// Закрываем тело HTTP-ответа
			require.NoError(t, resp.Body.Close())
		})
	}
}

func TestTextPostHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name                   string
		body                   string
		expectedURL            string
		expectedShortURLLength int
	}{
		{
			name:        "ValidRequest",
			body:        "http://example.com",
			expectedURL: "http://localhost:8080/",
		},
		{
			name:        "LongURL",
			body:        "https://www.thisisaverylongurlthatexceedsthecharacterlimit.com",
			expectedURL: "http://localhost:8080/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl, db, _ := Ctrl(t)
			db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).Times(1)
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос с текстовым телом и Content-Type: text/plain
			r := httptest.NewRequest("POST", "/", reqBody)
			r.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(ctrl.Save)
			handler(w, r)
			// Получаем результат (HTTP-ответ) после выполнения запроса.
			resp := w.Result()
			defer resp.Body.Close()

			// Проверяем статус
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	}
}

func TestShorten(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name         string
		body         string
		expectedJSON string
	}{
		{
			name:         "ValidRequest",
			body:         `{"url": "https://example.com", "alias": "example"}`,
			expectedJSON: `{"result":"/example"}`,
		},
		{
			name:         "Empty Request Body1",
			body:         `{}`,
			expectedJSON: `{"error":"поле URL обязательно для заполнения", "result":""}`,
		},
		{
			name:         "Empty Request Body2",
			body:         `{"url": "example.com"}`,
			expectedJSON: `{"error":"поле URL не является допустимым URL", "result":""}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, db, _ := Ctrl(t)
			// Если валидация завершается с ошибкой, то вызов Put не должен произойти
			if !strings.Contains(test.name, "ValidRequest") {
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				// Ожидаемый вызов Put
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any()).Return("example", nil).Times(1)
			}
			log := zap.NewNop()
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос
			req := httptest.NewRequest("POST", "/api/shorten", reqBody)
			req.Header.Set("Content-Type", "application/json")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()
			// Вызываем обработчик для HTTP POST запроса
			New(db, log, config.Cfg).Shorten(w, req)
			// Получаем результат (HTTP-ответ) после выполнения запроса.
			resp := w.Result()
			defer resp.Body.Close()

			// Чтобы сравнить JSON, сначала декодируем его из ответа сервера.
			var responseJSON map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&responseJSON)
			require.NoError(t, err)

			// Затем парсим ожидаемый JSON-ответ для сравнения.
			var expectedJSON map[string]interface{}
			err = json.NewDecoder(strings.NewReader(test.expectedJSON)).Decode(&expectedJSON)
			require.NoError(t, err)

			// Проверка, что полученный JSON совпадает с ожидаемым.
			assert.Equal(t, expectedJSON, responseJSON)
		})
	}
}
