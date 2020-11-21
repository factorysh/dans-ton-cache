package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

type DiskCache struct {
	path  string
	cache *lru.Cache
	lock  *sync.Mutex
	todos map[string]*sync.WaitGroup
}

// Get key, returns header, body io.ReadCloser, error
func (d *DiskCache) Get(key string) (http.Header, io.ReadCloser, error) {
	header, ok := d.cache.Get(key)
	if !ok {
		d.lock.Lock()
		todo, ok := d.todos[key]
		d.lock.Unlock()
		if !ok {
			return nil, nil, nil
		}
		todo.Wait()
		_, ok = d.cache.Get(key)
		if !ok {
			return nil, nil, fmt.Errorf("I wait for %v but it doesn't work", key)
		}
	}
	f, err := os.Open(path.Join(d.path, key))
	if err != nil {
		return nil, nil, err
	}
	return header.(http.Header), f, nil
}

type addCloser struct {
	wait    *sync.WaitGroup
	insider io.WriteCloser
}

func (a *addCloser) Write(data []byte) (int, error) {
	return a.insider.Write(data)
}

func (a *addCloser) Close() error {
	a.wait.Done()
	return a.insider.Close()
}

// Add key,header, returns io.WriteCloser and error
func (d *DiskCache) Add(key string, header http.Header) (io.WriteCloser, error) {
	todo := &sync.WaitGroup{}
	todo.Add(1)
	d.lock.Lock()
	d.todos[key] = todo
	d.lock.Unlock()
	d.cache.Add(key, header) // we don't care about eviction, yolo
	f, err := os.OpenFile(path.Join(d.path, fmt.Sprintf("%s.header", key)), os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	err = header.Write(f)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	// TODO write a temp file, when closed, mv to the file
	f, err = os.OpenFile(path.Join(d.path, key), os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	return &addCloser{
		wait:    todo,
		insider: f,
	}, nil
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
		path:  path,
		lock:  &sync.Mutex{},
		todos: make(map[string]*sync.WaitGroup),
	}
	l, err := lru.NewWithEvict(size, d.evict)
	if err != nil {
		return nil, err
	}
	d.cache = l
	return d, nil
}
