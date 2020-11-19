package cache

import (
	"io"
	"net/http"
)

// CacheHTTPWriter implements http.ResponseWriter
// TODO handle not 200 status
type CacheHTTPWriter struct {
	writer io.Writer
	ok     bool
	header http.Header
}

func (c *CacheHTTPWriter) Header() http.Header {
	if c.header == nil {
		c.header = make(http.Header)
	}
	return c.header
}

func (c *CacheHTTPWriter) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}

func (c *CacheHTTPWriter) WriteHeader(statusCode int) {
	c.ok = statusCode == http.StatusOK
}
