package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func GenerateCloudContainer(ctx context.Context) (
	string, int, int, *client.Client, error,
) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return "", 0, 0, nil, err
	}

	jupyterPort, sshPort := new(int), new(int)
	allocatePorts(jupyterPort, sshPort)

	containerConfig := container.Config{
		Image: "cloud",
		ExposedPorts: nat.PortSet{
			"8888/tcp": struct{}{},
			"22/tcp":   struct{}{},
		},
	}

	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: fmt.Sprintf("%d", *jupyterPort),
				},
			},
			"22/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: fmt.Sprintf("%d", *sshPort),
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&containerConfig,
		&hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", 0, 0, nil, err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", 0, 0, nil, err
	}

	return resp.ID, *jupyterPort, *sshPort, cli, nil
}
