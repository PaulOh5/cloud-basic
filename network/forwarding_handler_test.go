package network

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForwardingHandler(t *testing.T) {
	jupyterServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm jupyter-notebook"))
	}))
	defer jupyterServer.Close()

	vscodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm vscode"))
	}))

	fh := NewForwardingHandler()
	err := fh.AddForwarding("/jupyter", jupyterServer.URL)
	assert.NoError(t, err)
	err = fh.AddForwarding("/vscode", vscodeServer.URL)
	assert.NoError(t, err)

	mainServer := httptest.NewServer(fh)
	defer mainServer.Close()

	resp, err := http.Get(mainServer.URL + "/jupyter")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, I'm jupyter-notebook", getResponseBody(resp))

	resp, err = http.Get(mainServer.URL + "/vscode")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, I'm vscode", getResponseBody(resp))

	resp, err = http.Get(mainServer.URL + "/unknown")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestForwardingHandlerWithRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/tree", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, I'm jupyter-notebook"))
	}))
	mux.Handle("/", http.RedirectHandler("/tree", http.StatusFound))

	jupyterServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))
	defer jupyterServer.Close()

	fh := NewForwardingHandler()
	err := fh.AddForwarding("/jupyter", jupyterServer.URL)
	assert.NoError(t, err)

	mainServer := httptest.NewServer(fh)
	defer mainServer.Close()

	resp, err := http.Get(mainServer.URL + "/jupyter")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, I'm jupyter-notebook", getResponseBody(resp))
}
