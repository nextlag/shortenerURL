// Package gzip provides middleware and utility functions for handling gzip compression in HTTP responses and requests.
package gzip

import (
	"net/http"
	"strings"
)

// New creates a middleware handler for gzip compression and decompression.
func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w // Initialize a variable ow with the incoming http.ResponseWriter.

			// Get the value of the "Accept-Encoding" header from the request.
			acceptEncoding := r.Header.Get("Accept-Encoding")

			// Check for gzip support. If "Accept-Encoding" contains "gzip", set the supportGzip flag.
			supportGzip := strings.Contains(acceptEncoding, "gzip")

			// If gzip support is detected, create a new gzip.Writer (cw) and set ow to it.
			if supportGzip {
				cw := NewCompressWriter(w)
				ow = cw
				defer cw.Close() // Defer closing the gzip.Writer after processing.
			}

			// Get the value of the "Content-Encoding" header from the request.
			contentEncoding := r.Header.Get("Content-Encoding")

			// Check if the content was sent using gzip compression. If "Content-Encoding" contains "gzip", set the sendGzip flag.
			sendGzip := strings.Contains(contentEncoding, "gzip")

			// If the content was sent using gzip, create a new gzip.Reader (cr) and set r.Body to it.
			if sendGzip {
				cr, err := NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close() // Defer closing the gzip.Reader after processing.
			}

			// Call the original handler (next) with the modified http.ResponseWriter (ow) and the original request (r).
			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
