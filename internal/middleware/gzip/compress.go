// Package gzip provides utilities for compressing HTTP responses
// and decompressing HTTP requests using gzip encoding.
package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressWriter wraps an http.ResponseWriter and a gzip.Writer to
// compress HTTP responses. It implements the http.ResponseWriter interface.
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter creates a new CompressWriter that wraps the given
// http.ResponseWriter. It initializes a new gzip.Writer to handle
// the compression.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the HTTP headers of the wrapped response writer.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write compresses the data and writes it to the underlying gzip.Writer.
func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets the HTTP status code for the response. If the status code
// is less than 300, it sets the Content-Encoding header to "gzip".
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip.Writer to ensure all data is flushed and written
// to the underlying response writer.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}

// CompressReader wraps an io.ReadCloser and a gzip.Reader to
// decompress HTTP requests. It implements the io.ReadCloser interface.
type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader creates a new CompressReader that wraps the given
// io.ReadCloser. It initializes a new gzip.Reader to handle the decompression.
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read decompresses the data and reads it from the underlying gzip.Reader.
func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes both the underlying io.ReadCloser and the gzip.Reader.
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
