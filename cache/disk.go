package cache

import (
	"fmt"
	"io"
	"os"
	"path"

	lru "github.com/hashicorp/golang-lru"
)

type DiskCache struct {
	path  string
	cache *lru.Cache
}

func (d *DiskCache) Get(key string) (io.ReadCloser, error) {
	_, ok := d.cache.Get(key)
	if !ok {
		return nil, nil
	}
	f, err := os.Open(path.Join(d.path, key))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *DiskCache) Add(key string) (io.WriteCloser, error) {
	d.cache.Add(key, key) // we don't care about eviction, yolo
	// TODO write a temp file, when closed, mv to the file
	return os.OpenFile(path.Join(d.path, key), os.O_CREATE|os.O_WRONLY, 0640)
}

func (d *DiskCache) evict(key interface{}, value interface{}) {
	k, ok := key.(string)
	if !ok {
		panic(fmt.Sprintf("evict: Wrong key type : %v", key))
	}
	err := os.Remove(path.Join(d.path, string(k)))
	if err != nil {
		panic(err)
	}
}

func newDiskCache(path string, size int) (*DiskCache, error) {
	d := &DiskCache{
		path: path,
	}
	l, err := lru.NewWithEvict(size, d.evict)
	if err != nil {
		return nil, err
	}
	d.cache = l
	return d, nil
}
