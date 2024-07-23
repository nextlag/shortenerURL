package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
)

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
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				// Ожидаемый вызов Put
				db.EXPECT().DoPut(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("example", nil).Times(1)
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
