package network

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForwardingHandler(t *testing.T) {
	jupyterServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm jupyter-notebook"))
	}))
	defer jupyterServer.Close()

	vscodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm vscode"))
	}))

	forwardingList := ForwardingList{
		"jupyter-notebook": jupyterServer.URL,
		"vscode":           vscodeServer.URL,
	}

	handler, err := NewForwardingHandler(forwardingList)
	require.NoError(t, err)

	mainServer := httptest.NewServer(handler)
	defer mainServer.Close()

	resp, err := http.Get(mainServer.URL + "/jupyter-notebook")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello, I'm jupyter-notebook", getResponseBody(resp))

	resp, err = http.Get(mainServer.URL + "/vscode")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "Hello, I'm vscode", getResponseBody(resp))

	resp, err = http.Get(mainServer.URL + "/unknown")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
