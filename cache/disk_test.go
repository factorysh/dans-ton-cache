package cache

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisk(t *testing.T) {
	d, err := newDiskCache("/tmp", 3)
	assert.NoError(t, err)
	rc, err := d.Get("beuha")
	assert.NoError(t, err)
	assert.Nil(t, rc)

	err = d.Add("beuha", bytes.NewBuffer([]byte("aussi")))
	assert.NoError(t, err)

	buff := &bytes.Buffer{}
	rc, err = d.Get("beuha")
	assert.NoError(t, err)
	assert.NotNil(t, rc)
	defer rc.Close()
	_, err = io.Copy(buff, rc)
	assert.NoError(t, err)
	assert.Equal(t, []byte("aussi"), buff.Bytes())

	for k, v := range map[string]string{
		"a": "anachronic",
		"b": "beer",
		"c": "cat",
		"d": "data",
	} {
		d.Add(k, bytes.NewBuffer([]byte(v)))
	}

	assert.Equal(t, 3, d.cache.Len())
	keys := d.cache.Keys()
	fmt.Println(keys)
	for _, k := range []string{"b", "c", "d"} {
		assert.Contains(t, keys, k)
	}
}
