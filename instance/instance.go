package instance

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ContainerConfig struct {
	jupyterPort string
	sshPort     string
}

func createContainer(
	cli *client.Client,
	config ContainerConfig,
) (string, error) {
	containerConfig := container.Config{
		Image: "paul",
		ExposedPorts: nat.PortSet{
			"8888/tcp": struct{}{},
		},
	}
	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			"8888/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: config.jupyterPort,
				},
			},
			"22/tcp": []nat.PortBinding{
				{
					HostIP:   "",
					HostPort: config.sshPort,
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
