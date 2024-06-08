package network

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverserProxyHandler(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm jupyter-notebook"))
	}))
	defer targetServer.Close()

	handler, err := NewReverseProxyHandler(targetServer.URL)
	require.NoError(t, err)

	mainServer := httptest.NewServer(handler)
	defer mainServer.Close()

	resp, err := http.Get(mainServer.URL)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello, I'm jupyter-notebook", getResponseBody(resp))
}
