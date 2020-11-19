package cache

import (
	"fmt"
	"io"
	"net/http"

	"crypto/sha256"
)

type Cache struct {
	store *DiskCache
}

func New(path string, size int) (*Cache, error) {
	d, err := newDiskCache(path, size)
	if err != nil {
		return nil, err
	}
	return &Cache{d}, nil
}

func (c *Cache) key(r *http.Request) string {
	h := sha256.New()
	h.Write([]byte(r.URL.Path))

	return string(h.Sum(nil))
}

func (c *Cache) Middleware(in http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := c.key(r)
		rc, err := c.store.Get(key)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
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
		// TODO lock key, other call with same url waits
		rr, ww := io.Pipe()
		cw := &CacheHTTPWriter{
			writer: ww,
		}
		c.store.Add(key, rr)
		in(cw, r)
	}
}
