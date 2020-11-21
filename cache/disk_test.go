package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisk(t *testing.T) {
	d, err := newDiskCache("/tmp", 3)
	assert.NoError(t, err)
	header, rc, err := d.Get("beuha")
	assert.NoError(t, err)
	assert.Nil(t, rc)
	assert.Nil(t, header)

	wc, err := d.Add("beuha", nil)
	assert.NoError(t, err)
	io.WriteString(wc, "aussi")
	wc.Close()

	header, rc, err = d.Get("beuha")
	assert.NoError(t, err)
	assert.NotNil(t, rc)
	assert.Nil(t, header)
	defer rc.Close()
	r, err := ioutil.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("aussi"), r)
}

func TestFromPath(t *testing.T) {
	path, err := ioutil.TempDir("/tmp", "*")
	defer os.Remove(path)
	assert.NoError(t, err)
	d, err := newDiskCache(path, 3)
	assert.NoError(t, err)
	header := make(http.Header)
	header.Add("Name", "Bob")
	wc, err := d.Add("beuha", header)
	assert.NoError(t, err)
	io.WriteString(wc, "aussi")
	wc.Close()
	fmt.Println(path)

	d2, err := DiskCacheFromPath(path, 3)
	assert.NoError(t, err)
	h2, rc, err := d2.Get("beuha")
	assert.NoError(t, err)
	defer rc.Close()
	assert.Equal(t, "Bob", h2.Get("name"))
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
		wc, err := d.Add(k, nil)
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
	wc, err := d.Add("plop", nil)
	assert.NoError(t, err)
	w := &sync.WaitGroup{}
	w.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			header, rc, err := d.Get("plop")
			assert.NoError(t, err)
			defer rc.Close()
			assert.Nil(t, header)
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
