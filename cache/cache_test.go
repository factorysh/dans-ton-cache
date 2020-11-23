package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SlowHandler struct {
	n int
}

func (s *SlowHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-beuha", "aussi")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, "Hello world")
	if err != nil {
		fmt.Println(err)
	}
	s.n++
}

func TestMiddleware(t *testing.T) {
	c, err := New("/tmp", 3)
	assert.NoError(t, err)
	s := &SlowHandler{}
	h := c.Middleware(s.Handle)
	r := httptest.NewRequest("GET", "/beuha", nil)
	w := httptest.NewRecorder()
	h(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 1, s.n)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []byte("Hello world"), body)
	assert.Equal(t, "aussi", resp.Header.Get("x-beuha"))
	assert.Equal(t, "miss", resp.Header.Get("x-cache"))

	// Cached
	w = httptest.NewRecorder()
	h(w, r)
	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)
	assert.Equal(t, 1, s.n)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, []byte("Hello world"), body)
	assert.Equal(t, "aussi", resp.Header.Get("x-beuha"))
	assert.Equal(t, "hit", resp.Header.Get("x-cache"))
}
