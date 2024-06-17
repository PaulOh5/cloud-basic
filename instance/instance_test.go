package instance

import (
	"context"
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

	stdout, err := inst.Exec(context.Background(), "echo", "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)

	stdout, err = inst.Exec(context.Background(), "ls", "/not/exist")
	assert.Error(t, err)
	assert.Empty(t, stdout)
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
