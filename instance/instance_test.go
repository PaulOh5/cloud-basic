package instance

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestInstanceStatus(t *testing.T) {
	inst, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		inst.Remove(context.Background())
	})

	status, err := inst.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, RUNNING, status)

	err = inst.Stop(context.Background())
	require.NoError(t, err)

	status, err = inst.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, STOPPED, status)

	err = inst.Start(context.Background())
	require.NoError(t, err)

	status, err = inst.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, RUNNING, status)

	err = inst.Remove(context.Background())
	require.NoError(t, err)

	_, err = inst.cli.ContainerInspect(context.Background(), inst.containerID)
	require.Error(t, err)
}

func TestInstanceExecCommand(t *testing.T) {
	inst, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := inst.Remove(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	stdout, stderr, err := inst.Exec(context.Background(), "echo", "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
	assert.Empty(t, stderr)

	stdout, stderr, err = inst.Exec(context.Background(), "ls", "/not/exist")
	require.NoError(t, err)
	assert.Empty(t, stdout)
	assert.Equal(t, "ls: cannot access '/not/exist': No such file or directory\n", stderr)
}

func TestJupyterConnection(t *testing.T) {
	inst, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := inst.Remove(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	resp, err := http.Get(inst.GetJupyterUrl())
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSshConnection(t *testing.T) {
	inst, err := NewInstance(context.Background())
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := inst.Remove(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	config, err := inst.GetSshConfig()
	require.NoError(t, err)

	sshClient, err := ssh.Dial("tcp", inst.GetSshUrl(), config)
	require.NoError(t, err)
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	require.NoError(t, err)
	defer session.Close()

	output, err := session.CombinedOutput("echo hello")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", string(output))
}
