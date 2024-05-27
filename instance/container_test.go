package instance

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestContainer(t *testing.T) {
	testCases := []struct {
		config ContainerConfig
	}{
		{config: ContainerConfig{jupyterPort: "9999", sshPort: "2222"}},
	}

	for i, tc := range testCases {
		t.Run(tc.config.jupyterPort, func(t *testing.T) {
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				t.Fatal(err)
			}
			defer cli.Close()
			containerID, err := createContainer(cli, tc.config)
			if err != nil {
				t.Error(err)
			}
			defer removeContainer(cli, containerID)
			time.Sleep(5 * time.Second)

			if err := checkContainerCreated(cli, containerID); err != nil {
				t.Fatalf("container not created: %v", err)
			}
			if err := checkJupyterResponse(tc.config.jupyterPort); err != nil {
				t.Fatalf("failed to connect to jupyter notebook: %v", err)
			}
			t.Logf("test case %d passed", i+1)
		})
	}
}

func checkContainerCreated(cli *client.Client, containerID string) error {
	ctr, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return err
	}

	if ctr.ID != containerID {
		return fmt.Errorf("container not found")
	}

	if !ctr.State.Running {
		return fmt.Errorf("container not running")
	}
	return nil
}

func checkJupyterResponse(port string) error {
	resp, err := http.Get("http://localhost:" + port)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jupyter notebook not running")
	}
	return nil
}

// func checkVolumeBinding(cli *client.Client, containerID string) error {
// 	ctr, err := cli.ContainerInspect(context.Background(), containerID)
// 	if err != nil {
// 		return err
// 	}

// 	if len(ctr.Mounts) == 0 {
// 		return fmt.Errorf("volume not bound")
// 	}

// 	resp, err := cli.ContainerExecCreate(context.Background(), ctr.ID, types.ExecConfig{
// 		Cmd: []string{"mkdir", ""},
// 	})
// 	return nil
// }

func removeContainer(cli *client.Client, containerID string) error {
	err := cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}
