package instance

import (
	"context"
	"fmt"

	sshkey "github.com/PaulOh5/cloud-basic/ssh_key"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Instance struct {
	cli                  *client.Client
	sshKey               *sshkey.SSHKey
	jupyterPort, sshPort int
	containerID          string
}

func (i *Instance) Exec(ctx context.Context) error {
	return nil
}

func (i *Instance) Start(ctx context.Context) error {
	err := i.cli.ContainerStart(ctx, i.containerID, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Stop(ctx context.Context) error {
	err := i.cli.ContainerStop(ctx, i.containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Remove(ctx context.Context) error {
	err := i.cli.ContainerRemove(context.Background(), i.containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) GetStatus(ctx context.Context) (InstanceStatus, error) {
	ctr, err := i.cli.ContainerInspect(ctx, i.containerID)
	if err != nil {
		return STOPPED, err
	}

	if ctr.State.Running {
		return RUNNING, nil
	}

	return STOPPED, nil
}

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

// TODO: 리팩토링 필요
func NewInstance(ctx context.Context) (*Instance, error) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	sshKey, err := sshkey.NewSshKey()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return nil, err
	}

	instance := &Instance{
		cli:         cli,
		sshKey:      sshKey,
		jupyterPort: *jupyterPort,
		sshPort:     *sshPort,
		containerID: resp.ID,
	}

	return instance, nil
}

type ContainerConfig struct {
	JupyterPort string
	SshPort     string
}

func applyKeyToContainer(
	cli *client.Client,
	containerID string,
	privateKey, publicKey []byte,
) error {
	return nil
}

func RemoveContainer(cli *client.Client, containerID string) error {
	err := cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}
