package instance

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstanceStatus(t *testing.T) {
	instance, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Remove(context.Background())

	status, err := instance.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, RUNNING, status)

	err = instance.Stop(context.Background())
	require.NoError(t, err)

	status, err = instance.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, STOPPED, status)

	err = instance.Start(context.Background())
	require.NoError(t, err)

	status, err = instance.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, RUNNING, status)

	err = instance.Remove(context.Background())
	require.NoError(t, err)

	_, err = instance.cli.ContainerInspect(context.Background(), instance.containerID)
	require.Error(t, err)
}

func TestInstanceExec(t *testing.T) {
	instance, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer instance.Remove(context.Background())

	stdout, stderr, err := instance.Exec(context.Background(), "echo", "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
	assert.Empty(t, stderr)

	stdout, stderr, err = instance.Exec(context.Background(), "ls", "/not/exist")
	require.NoError(t, err)
	assert.Empty(t, stdout)
	assert.Equal(t, "ls: cannot access '/not/exist': No such file or directory\n", stderr)
}

// func TestContainer(t *testing.T) {
// 	testCases := []struct {
// 		config ContainerConfig
// 		key    *sshkey.SSHKey
// 	}{
// 		{config: ContainerConfig{
// 			JupyterPort: "9999",
// 			SshPort:     "2222",
// 		}},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(tc.config.JupyterPort, func(t *testing.T) {
// 			tc.key, _ = sshkey.NewSshKey()

// 			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer cli.Close()
// 			containerID, err := CreateContainer(cli, tc.config)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			defer RemoveContainer(cli, containerID)
// 			time.Sleep(5 * time.Second)

// 			if err := checkContainerCreated(cli, containerID); err != nil {
// 				t.Fatalf("container not created: %v", err)
// 			}
// 			if err := checkJupyterResponse(tc.config.JupyterPort); err != nil {
// 				t.Fatalf("failed to connect to jupyter notebook: %v", err)
// 			}
// 			if err := checkSshConnection(tc.config); err != nil {
// 				t.Fatalf("failed to connect to ssh: %v", err)
// 			}
// 			t.Logf("test case %d passed", i+1)
// 		})
// 	}
// }

// func checkContainerCreated(cli *client.Client, containerID string) error {
// 	ctr, err := cli.ContainerInspect(context.Background(), containerID)
// 	if err != nil {
// 		return err
// 	}

// 	if ctr.ID != containerID {
// 		return fmt.Errorf("container not found")
// 	}

// 	if !ctr.State.Running {
// 		return fmt.Errorf("container not running")
// 	}
// 	return nil
// }

// func checkJupyterResponse(port string) error {
// 	resp, err := http.Get("http://localhost:" + port)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("jupyter notebook not running")
// 	}
// 	return nil
// }
