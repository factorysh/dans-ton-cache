package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisk(t *testing.T) {
	d, err := newDiskCache("/tmp", 3)
	assert.NoError(t, err)
	rc, err := d.Get("beuha")
	assert.NoError(t, err)
	assert.Nil(t, rc)

	wc, err := d.Add("beuha")
	assert.NoError(t, err)
	io.WriteString(wc, "aussi")
	wc.Close()

	rc, err = d.Get("beuha")
	assert.NoError(t, err)
	assert.NotNil(t, rc)
	defer rc.Close()
	r, err := ioutil.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("aussi"), r)
}

func TestEviction(t *testing.T) {
	d, err := newDiskCache("/tmp", 3)
	assert.NoError(t, err)
	datas := map[string]string{
		"a": "anachronic",
		"b": "beer",
		"c": "cat",
		"d": "data",
	}
	var keys []string
	for k := range datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k)
		wc, err := d.Add(k)
		assert.NoError(t, err)
		defer wc.Close()
		io.WriteString(wc, datas[k])
	}

	assert.Equal(t, 3, d.cache.Len())
	myKeys := d.cache.Keys()
	fmt.Println(keys)
	for _, k := range []string{"b", "c", "d"} {
		assert.Contains(t, myKeys, k)
	}
}

func TestConcurrentAcces(t *testing.T) {
	d, err := newDiskCache("/tmp", 3)
	assert.NoError(t, err)
	wc, err := d.Add("plop")
	assert.NoError(t, err)
	w := &sync.WaitGroup{}
	w.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			rc, err := d.Get("plop")
			assert.NoError(t, err)
			defer rc.Close()
			v, err := ioutil.ReadAll(rc)
			assert.NoError(t, err)
			assert.Equal(t, []byte("aussi"), v)
			w.Done()
		}()
	}
	_, err = io.WriteString(wc, "aussi")
	assert.NoError(t, err)
	wc.Close()
	w.Wait()
}
