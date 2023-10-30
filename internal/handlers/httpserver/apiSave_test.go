package httpserver_test

import (
	"encoding/json"
	"github.com/nextlag/shortenerURL/internal/handlers/httpserver"
	"github.com/nextlag/shortenerURL/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShorten(t *testing.T) {
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
			expectedJSON: `{"error":"поле URL обязательно для заполнения"}`,
		},
		{
			name:         "Empty Request Body2",
			body:         `{"url": "example.com"}`,
			expectedJSON: `{"error":"поле URL не является допустимым URL"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Создаем фейковое хранилище
			db := storage.New()
			log := zap.NewNop()
			// Создаем объект reqBody, который реализует интерфейс io.Reader и будет представлять тело запроса.
			reqBody := strings.NewReader(test.body)
			// Создаем новый POST запрос
			req := httptest.NewRequest("POST", "/api/shorten", reqBody)
			req.Header.Set("Content-Type", "application/json")
			// Создаем записывающий ResponseRecorder, который будет использоваться для записи HTTP ответа.
			w := httptest.NewRecorder()
			// Вызываем обработчик для HTTP POST запроса
			httpserver.Shorten(log, db).ServeHTTP(w, req)
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
