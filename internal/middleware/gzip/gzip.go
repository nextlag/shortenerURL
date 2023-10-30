package gzip

import (
	"net/http"
	"strings"

	"github.com/nextlag/shortenerURL/internal/lib/gzip"
)

func NewGzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w // Создаем переменную ow и инициализируем ее с http.ResponseWriter из входящих параметров.

		// Получаем значение заголовка "Accept-Encoding" из запроса.
		acceptEncoding := r.Header.Get("Accept-Encoding")

		// Проверяем поддержку сжатия gzip. Если заголовок "Accept-Encoding" содержит "gzip", устанавливаем флаг supportGzip.
		supportGzip := strings.Contains(acceptEncoding, "gzip")

		// Если поддержка gzip обнаружена, создаем новый gzip.Writer (cw) и устанавливаем ow на него.
		if supportGzip {
			cw := gzip.NewCompressWriter(w)
			ow = cw
			defer cw.Close() // Отложенное закрытие gzip.Writer после завершения обработки.
		}

		// Получаем значение заголовка "Content-Encoding" из запроса.
		contentEncoding := r.Header.Get("Content-Encoding")

		// Проверяем, был ли отправлен контент с использованием сжатия gzip. Если "Content-Encoding" содержит "gzip", устанавливаем флаг sendGzip.
		sendGzip := strings.Contains(contentEncoding, "gzip")

		// Если контент был отправлен с использованием gzip, создаем новый gzip.Reader (cr) и устанавливаем r.Body на него.
		if sendGzip {
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close() // Отложенное закрытие gzip.Reader после завершения обработки.
		}

		// Вызываем оригинальный обработчик (h) с модифицированным http.ResponseWriter (ow) и исходным запросом (r).
		h.ServeHTTP(ow, r)
	}
}
