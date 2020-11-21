package cache

import (
	"fmt"
	"io"
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

func (d *DiskCache) Get(key string) (io.ReadCloser, error) {
	_, ok := d.cache.Get(key)
	if !ok {
		d.lock.Lock()
		todo, ok := d.todos[key]
		d.lock.Unlock()
		if !ok {
			return nil, nil
		}
		todo.Wait()
		_, ok = d.cache.Get(key)
		if !ok {
			return nil, fmt.Errorf("I wait for %v but it doesn't work", key)
		}
	}
	f, err := os.Open(path.Join(d.path, key))
	if err != nil {
		return nil, err
	}
	return f, nil
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

func (d *DiskCache) Add(key string) (io.WriteCloser, error) {
	todo := &sync.WaitGroup{}
	todo.Add(1)
	d.lock.Lock()
	d.todos[key] = todo
	d.lock.Unlock()
	d.cache.Add(key, key) // we don't care about eviction, yolo
	// TODO write a temp file, when closed, mv to the file
	f, err := os.OpenFile(path.Join(d.path, key), os.O_CREATE|os.O_WRONLY, 0640)
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
