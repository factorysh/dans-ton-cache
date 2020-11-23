package cache

import (
	"fmt"
	"io"
	"net/http"
)

// CacheHTTPWriter implements http.ResponseWriter
// TODO handle not 200 status
type CacheHTTPWriter struct {
	writer     io.Writer
	httpWriter http.ResponseWriter
	ok         bool
	onBody     func(status int, header http.Header) (io.WriteCloser, error)
	header     http.Header
	closer     io.Closer
}

func NewCacheHTTPWriter(httpWriter http.ResponseWriter,
	onBody func(status int, header http.Header) (io.WriteCloser, error)) *CacheHTTPWriter {
	return &CacheHTTPWriter{
		header:     make(http.Header),
		onBody:     onBody,
		httpWriter: httpWriter,
	}
}

func (c *CacheHTTPWriter) Header() http.Header {
	return c.header
}

func (c *CacheHTTPWriter) Write(data []byte) (int, error) {
	if c.ok {
		return c.writer.Write(data)
	}
	fmt.Println("Reading body, but it's not a 200")
	return c.httpWriter.Write(data)
}

func (c *CacheHTTPWriter) WriteHeader(statusCode int) {
	c.ok = statusCode == http.StatusOK
	w, err := c.onBody(statusCode, c.header)
	if err != nil {
		return
	}
	for key, values := range c.header {
		for _, value := range values {
			c.httpWriter.Header().Add(key, value)

		}
	}
	c.httpWriter.WriteHeader(statusCode)
	c.writer = io.MultiWriter(w, c.httpWriter)
	c.closer = w
}

func (c *CacheHTTPWriter) Close() error {
	return c.closer.Close()
}
