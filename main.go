package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func createJupyterNotebookContainer(cli *client.Client) (string, error) {
	containerConfig := container.Config{
		Image: "jupyter/datascience-notebook",
		ExposedPorts: nat.PortSet{
			"8888/tcp": struct{}{},
		},
	}
	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: "8888",
				},
			},
		},
	}
	resp, err := cli.ContainerCreate(
		context.Background(),
		&containerConfig,
		&hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	cli.NegotiateAPIVersion(context.Background())

	containerList, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containerList {
		fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
	}
}
