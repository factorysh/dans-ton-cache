package cache

import (
	"io"
	"os"
	"path"

	lru "github.com/hashicorp/golang-lru"
)

type DiskCache struct {
	path  string
	cache *lru.Cache
}

func (d *DiskCache) Get(key string, writer io.Writer) (bool, error) {
	_, ok := d.cache.Get(key)
	if !ok {
		return false, nil
	}
	f, err := os.Open(path.Join(d.path, key))
	if err != nil {
		return false, err
	}
	_, err = io.Copy(writer, f)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *DiskCache) Add(key string, reader io.Reader) error {
	evicted := d.cache.Add(key, key)
	if !evicted {
		f, err := os.OpenFile(path.Join(d.path, key), os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DiskCache) evict(key interface{}, value interface{}) {
	k, ok := key.([]byte)
	if !ok {
		panic("Wrong type")
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
