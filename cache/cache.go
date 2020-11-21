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
	d, err := newDiskCache(path, size)
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
		rc, err := c.store.Get(key)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		if rc != nil {
			defer rc.Close()
			_, err = io.Copy(w, rc)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
			}
			return
		}
		wc, err := c.store.Add(key)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
		defer wc.Close()
		in(&CacheHTTPWriter{
			writer: io.MultiWriter(w, wc),
		}, r)
	}
}
