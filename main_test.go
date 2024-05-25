package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestCreateJupyterNotebookContainer(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	containerID, err := createJupyterNotebookContainer(cli)
	if err != nil {
		t.Error(err)
	}
	defer removeContainer(cli, containerID)
	time.Sleep(10 * time.Second)

	if err := checkContainerCreated(cli, containerID); err != nil {
		t.Fatalf("container not created: %v", err)
	}
	t.Log("container created")

	if err := checkJupyterResponse(); err != nil {
		t.Fatalf("failed to connect to jupyter notebook: %v", err)
	}
	t.Log("jupyter notebook running")
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

func checkJupyterResponse() error {
	resp, err := http.Get("http://localhost:8888")
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
