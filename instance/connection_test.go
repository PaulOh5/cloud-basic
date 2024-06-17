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

func TestConnectJupyterInstance(t *testing.T) {
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

	resp, err := http.Get(server.URL + inst.jupyterURL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
