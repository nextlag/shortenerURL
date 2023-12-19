package main_test

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app/mocks"
	"github.com/nextlag/shortenerURL/internal/transport/rest/handlers"
	gz "github.com/nextlag/shortenerURL/internal/transport/rest/middleware/gzip"
)

func TestMain(m *testing.M) {
	if err := config.MakeConfig(); err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	// Закрываем все оставшиеся тела ответов
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()
	// Завершаем выполнение программы с кодом завершения.
	os.Exit(exitCode)
}

func TestGetHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создаем фейковое хранилище
	db := mocks.NewMockStorage(ctrl)

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

			if test.Name == "Valid ID" {
				db.EXPECT().Get(gomock.Any(), gomock.Any()).Return("", false, nil).Times(1)
			} else {
				db.EXPECT().Get(gomock.Any(), gomock.Any()).Return("", false, errors.New("error")).Times(1)
			}
			// Создаем фейковый запрос
			req := httptest.NewRequest("GET", test.RequestPath, nil)
			w := httptest.NewRecorder()

			// Создаем и вызываем handler для маршрута
			handlers.NewGetHandler(db, nil, config.Config).ServeHTTP(w, req)
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
			// Создаем фейковое хранилище
			db := mocks.NewMockStorage(ctrl)
			db.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil).Times(1)

			// Создаем фейковый логгер
			fakeLogger := zap.NewNop()

			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос с текстовым телом и Content-Type: text/plain
			req := httptest.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "text/plain")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()

			// Используем фейковый логгер при создании SaveHandler
			handlers.NewSaveHandlers(db, fakeLogger, config.Config).ServeHTTP(w, req)
			// Получаем результат (HTTP-ответ) после выполнения запроса.
			resp := w.Result()
			defer resp.Body.Close()

			// Проверяем статус
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	}
}
func TestGzipMiddleware(t *testing.T) {
	// Структура с параметрами для тестовых случаев.
	testCases := []struct {
		name              string
		acceptEncoding    string
		expectedBodyEmpty bool
		expectedHeader    string
	}{
		{
			name:           "Accepts_Gzip_Encoding",
			acceptEncoding: "gzip",
			// В этом случае ожидается, что тело ответа будет сжато.
			expectedBodyEmpty: false,
			expectedHeader:    "TestHeader",
		},
		{
			name:           "No_Gzip_Encoding",
			acceptEncoding: "",
			// В этом случае ожидается, что тело ответа будет пустым.
			expectedBodyEmpty: true,
			expectedHeader:    "TestHeader",
		},
	}

	// Создаем тестовый HTTP-обработчик для тестирования middleware.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "TestHeader")
		w.WriteHeader(http.StatusOK)
	})

	// Создаем экземпляр middleware с тестовым обработчиком.
	middleware := gz.New(testHandler)

	// Проходим по всем тестовым случаям.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Создаем тестовый запрос с указанным заголовком Accept-Encoding.
			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("Accept-Encoding", tc.acceptEncoding)

			// Создаем тестовую запись (Recorder) для записи ответа.
			rr := httptest.NewRecorder()

			// Запускаем middleware.
			middleware.ServeHTTP(rr, req)

			// Проверяем, что middleware корректно обработал запрос.
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
			}

			// Проверяем, что middleware вернул ожидаемое тело ответа.
			if (len(rr.Body.Bytes()) == 0) != tc.expectedBodyEmpty {
				t.Errorf("Expected empty response body: %v, but got: %v", tc.expectedBodyEmpty, len(rr.Body.Bytes()) == 0)
			}

			// Проверяем, что middleware добавил заголовок X-Custom-Header.
			if rr.Header().Get("X-Custom-Header") != tc.expectedHeader {
				t.Errorf("Expected X-Custom-Header to be '%s', but got %s", tc.expectedHeader, rr.Header().Get("X-Custom-Header"))
			}
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
			expectedJSON: `{"result":"http://localhost:8080/example"}`,
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
			// Создаем фейковое хранилище
			db := mocks.NewMockStorage(ctrl)
			// Если валидация завершается с ошибкой, то вызов Put не должен произойти
			if !strings.Contains(test.name, "ValidRequest") {
				db.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				// Ожидаемый вызов Put
				db.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return("example", nil).Times(1)
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
			handlers.NewShortenHandlers(db, log, config.Config).ServeHTTP(w, req)
			// Проверяем, что все ожидаемые вызовы были выполнены
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
