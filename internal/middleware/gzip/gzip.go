// Package gzip provides middleware for handling gzip compression
// and decompression in HTTP requests and responses.
package gzip

import (
	"net/http"
	"strings"
)

// New returns a middleware handler function that compresses HTTP responses
// and decompresses HTTP requests using gzip encoding.
func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w // Initialize ow with the incoming http.ResponseWriter.

			// Get the "Accept-Encoding" header from the request.
			acceptEncoding := r.Header.Get("Accept-Encoding")

			// Check for gzip support. If the "Accept-Encoding" header contains "gzip", set supportGzip to true.
			supportGzip := strings.Contains(acceptEncoding, "gzip")

			// If gzip support is detected, create a new gzip.Writer (cw) and set ow to it.
			if supportGzip {
				cw := NewCompressWriter(w)
				ow = cw
				defer cw.Close() // Defer closing the gzip.Writer until the handler is finished.
			}

			// Get the "Content-Encoding" header from the request.
			contentEncoding := r.Header.Get("Content-Encoding")

			// Check if the content was sent using gzip. If "Content-Encoding" contains "gzip", set sendGzip to true.
			sendGzip := strings.Contains(contentEncoding, "gzip")

			// If the content was sent using gzip, create a new gzip.Reader (cr) and set r.Body to it.
			if sendGzip {
				cr, err := NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close() // Defer closing the gzip.Reader until the handler is finished.
			}

			// Call the original handler (next) with the modified http.ResponseWriter (ow) and the original request (r).
			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
