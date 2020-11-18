package cache

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisk(t *testing.T) {
	d, err := newDiskCache("/tmp", 10)
	assert.NoError(t, err)
	buff := &bytes.Buffer{}
	ok, err := d.Get("beuha", buff)
	assert.NoError(t, err)
	assert.False(t, ok)

	err = d.Add("beuha", bytes.NewBuffer([]byte("aussi")))
	assert.NoError(t, err)

	buff = &bytes.Buffer{}
	ok, err = d.Get("beuha", buff)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []byte("aussi"), buff.Bytes())
}
