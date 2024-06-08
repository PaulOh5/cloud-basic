package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startHTTPServer(hs *HandlerSwitcher) *httptest.Server {
	s := httptest.NewServer(hs)
	return s
}

func getResponseBody(resp *http.Response) string {
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	resp.Body.Close()
	return string(output)
}

func TestHandlerSwitcher(t *testing.T) {
	hs := NewHandlerSwitcher()
	server := startHTTPServer(hs)
	defer server.Close()

	hs.AddHandler("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	resp, err := http.Get(server.URL + "/hello")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, world!", getResponseBody(resp))

	resp, err = http.Get(server.URL + "/goodbye")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	hs.AddHandler("/goodbye", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Goodbye, world!"))
	}))

	resp, err = http.Get(server.URL + "/goodbye")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Goodbye, world!", getResponseBody(resp))

	hs.RemoveHandler("/goodbye")

	resp, err = http.Get(server.URL + "/goodbye")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
