package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/usecase"
	"github.com/nextlag/shortenerURL/internal/usecase/auth"
)

// Save - обработчик POST-запросов для создания и сохранения URL в storage.
func (c *Controller) Save(w http.ResponseWriter, r *http.Request) {

	// Считываем тело запроса (оригинальный URL)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}
	uuid, err := auth.CheckCookie(w, r, c.log)
	if err != nil {
		c.log.Error("Error getting cookie: ", zap.Error(err))
		return
	}

	// Попытка сохранить short-URL и оригинальный URL в хранилище
	alias, err := c.uc.DoPut(r.Context(), string(body), uuid)

	// Обработка конфликта дубликатов
	if errors.Is(err, usecase.ErrConflict) {
		c.log.Error("duplicate url", zap.String("alias", alias), zap.String("url", string(body)))
		w.WriteHeader(http.StatusConflict)
		_, err = fmt.Fprintf(w, "%s/%s", c.cfg.BaseURL, alias)
		if err != nil {
			c.log.Error("error sending short URL response", zap.Error(err))
			return
		}
		return
	}

	// Обработка других ошибок
	if err != nil {
		c.log.Error("failed to add URL", zap.Error(err), zap.String("path to file storage", c.cfg.FileStorage))
		http.Error(w, fmt.Sprintf("failed to add URL: %s", err), http.StatusInternalServerError)
		return
	}
	// Устанавливаем статус HTTP 201 Created
	w.WriteHeader(http.StatusCreated)

	// Отправляем short-URL в теле HTTP-ответа
	_, err = fmt.Fprintf(w, "%s/%s", c.cfg.BaseURL, alias)
	if err != nil {
		c.log.Error("error sending short URL response", zap.Error(err))
		return
	}
	c.log.Info("alias added success", zap.String("alias", alias))
}
