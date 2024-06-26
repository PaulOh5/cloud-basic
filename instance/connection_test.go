package instance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaulOh5/cloud-basic/network"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestConnectSsh(t *testing.T) {
	inst, err := NewInstance(context.Background())
	assert.NoError(t, err)

	t.Cleanup(func() {
		if err := inst.Remove(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	config, err := inst.GetSshConfig()
	assert.NoError(t, err)

	sshClient, err := ssh.Dial("tcp", inst.GetSshUrl(), config)
	assert.NoError(t, err)
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	assert.NoError(t, err)
	defer session.Close()

	output, err := session.CombinedOutput("echo hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello\n", string(output))
}

func TestConnectJupyter(t *testing.T) {
	inst, err := NewInstance(context.Background())
	assert.NoError(t, err)

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

// func TestConnectSsh(t *testing.T) {
// 	inst, err := NewInstance(context.Background())
// 	assert.NoError(t, err)
// }
