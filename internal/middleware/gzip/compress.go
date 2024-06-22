// Package gzip provides middleware and utility functions for handling gzip compression in HTTP responses and requests.
package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressWriter wraps an http.ResponseWriter to provide gzip compression for the response body.
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter creates a new CompressWriter with the provided http.ResponseWriter.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by WriteHeader.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the data to the connection as part of an HTTP reply, compressing it using gzip.
func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code, setting Content-Encoding to gzip if the status code is less than 300.
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip.Writer, flushing any unwritten data to the underlying http.ResponseWriter.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}

// CompressReader wraps an io.ReadCloser to provide gzip decompression for the request body.
type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader creates a new CompressReader with the provided io.ReadCloser, initializing a new gzip.Reader.
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

// Read reads and decompresses data from the wrapped io.ReadCloser.
func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes both the gzip.Reader and the underlying io.ReadCloser.
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
