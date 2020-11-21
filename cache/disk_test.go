package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
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
