package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type Cache struct {
	store *DiskCache
}

func New(path string, size int) (*Cache, error) {
	d, err := DiskCacheFromPath(path, size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		store: d,
	}, nil
}

func (c *Cache) key(r *http.Request) string {
	h := sha256.New()
	io.WriteString(h, r.URL.Path)
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Cache) Middleware(in http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := c.key(r)
		header, rc, err := c.store.Get(key)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		if rc != nil {
			defer rc.Close()
			for k, values := range header {
				for _, value := range values {
					w.Header().Add(k, value)
				}
			}
			_, err = io.Copy(w, rc)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
			}
			return
		}
		header = w.Header()
		wc, err := c.store.Add(key, header)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
		defer wc.Close()
		in(&CacheHTTPWriter{
			header: header,
			writer: io.MultiWriter(w, wc),
		}, r)
	}
}
