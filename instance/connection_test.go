package instance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaulOh5/cloud-basic/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceConnectiont(t *testing.T) {
	inst, err := NewInstance(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := inst.Remove(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	fh := network.NewForwardingHandler()
	err = inst.EstablishConnect(fh)
	assert.NoError(t, err)
	defer inst.Disconnect(fh)

	server := httptest.NewServer(fh)
	defer server.Close()

	resp, err := http.Get(server.URL + "/" + inst.jupyterURL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// func TestJupyterForwarding(t *testing.T) {
// 	jupyterURL := "http://127.0.0.1:8888"
// 	handler, err := network.NewForwardingHandler(jupyterURL)
// 	assert.NoError(t, err)

// 	mux := http.NewServeMux()
// 	mux.Handle("/jupyter", handler)

// 	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		t.Logf("Request URL Path: %s\n", r.URL.Path)
// 		mux.ServeHTTP(w, r)
// 	})

// 	server := httptest.NewServer(newHandler)
// 	defer server.Close()

// 	resp, err := http.Get(server.URL + "/jupyter")
// 	assert.NoError(t, err)
// 	defer resp.Body.Close()
// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
// }
