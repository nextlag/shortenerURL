package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/nextlag/shortenerURL/internal/config"
	"github.com/nextlag/shortenerURL/internal/service/app"
	"github.com/nextlag/shortenerURL/internal/service/auth"
	"github.com/nextlag/shortenerURL/internal/storage/dbstorage"
)

type SaveHandler struct {
	db  app.Storage
	log *zap.Logger
	cfg config.Args
}

func NewSaveHandlers(db app.Storage, log *zap.Logger, cfg config.Args) *SaveHandler {
	return &SaveHandler{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

// Save - обработчик POST-запросов для создания и сохранения URL в storage.
func (s *SaveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Считываем тело запроса (оригинальный URL)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request 400", http.StatusBadRequest)
		return
	}
	uid := auth.CheckCookie(w, r, s.log)

	// Попытка сохранить short-URL и оригинальный URL в хранилище
	alias, err := s.db.Put(r.Context(), string(body), uid)

	// Обработка конфликта дубликатов
	if errors.Is(err, dbstorage.ErrConflict) {
		s.log.Error("duplicate url", zap.String("alias", alias), zap.String("url", string(body)))
		w.WriteHeader(http.StatusConflict)
		_, err = fmt.Fprintf(w, "%s/%s", s.cfg.URLShort, alias)
		if err != nil {
			s.log.Error("error sending short URL response", zap.Error(err))
			return
		}
		return
	}

	// Обработка других ошибок
	if err != nil {
		s.log.Error("failed to add URL", zap.Error(err), zap.String("path to filestorage", s.cfg.FileStorage))
		http.Error(w, fmt.Sprintf("failed to add URL: %s", err), http.StatusInternalServerError)
		return
	}
	// Устанавливаем статус HTTP 201 Created
	w.WriteHeader(http.StatusCreated)

	// Отправляем short-URL в теле HTTP-ответа
	_, err = fmt.Fprintf(w, "%s/%s", s.cfg.URLShort, alias)
	if err != nil {
		s.log.Error("error sending short URL response", zap.Error(err))
		return
	}
	s.log.Info("alias added success", zap.String("alias", alias))
}
