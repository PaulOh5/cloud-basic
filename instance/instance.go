package instance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PaulOh5/cloud-basic/network"
	sshkey "github.com/PaulOh5/cloud-basic/ssh_key"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"golang.org/x/crypto/ssh"
)

type Instance struct {
	cli                  *client.Client
	key                  *sshkey.SSHKey
	jupyterPort, sshPort int
	jupyterURL           string
	containerID          string
}

func (i Instance) Exec(ctx context.Context, cmd ...string) (string, string, error) {
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := i.cli.ContainerExecCreate(ctx, i.containerID, execConfig)
	if err != nil {
		return "", "", err
	}

	resp, err := i.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", "", err
	}
	defer resp.Close()

	var outBuf, errBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
	if err != nil {
		return "", "", err
	}

	return outBuf.String(), errBuf.String(), nil
}

func (i Instance) Start(ctx context.Context) error {
	err := i.cli.ContainerStart(ctx, i.containerID, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) Stop(ctx context.Context) error {
	err := i.cli.ContainerStop(ctx, i.containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (i Instance) Remove(ctx context.Context) error {
	err := i.cli.ContainerRemove(context.Background(), i.containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}

	i.cli.Close()
	return nil
}

func (i Instance) GetStatus(ctx context.Context) (InstanceStatus, error) {
	ctr, err := i.cli.ContainerInspect(ctx, i.containerID)
	if err != nil {
		return STOPPED, err
	}

	if ctr.State.Running {
		return RUNNING, nil
	}

	return STOPPED, nil
}

func (i *Instance) EstablishConnect(fh *network.ForwardingHandler) error {
	jupyterPath := fmt.Sprintf("http://127.0.0.1:%d", i.jupyterPort)
	err := fh.AddForwarding("/"+i.containerID[12:]+"/jupyter", jupyterPath)
	if err != nil {
		return err
	}

	i.jupyterURL = "/" + i.containerID[12:] + "/jupyter"
	return nil
}

func (i Instance) Disconnect(fh *network.ForwardingHandler) {
	fh.RemoveForwarding(i.jupyterURL)
}

func (i Instance) GetSshUrl() string {
	return fmt.Sprintf("localhost:%d", i.sshPort)
}

func (i Instance) GetSshConfig() (*ssh.ClientConfig, error) {
	signer, err := ssh.NewSignerFromKey(i.key.PrivateKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config, nil
}

// TODO: 리팩토링 필요
func NewInstance(ctx context.Context) (*Instance, error) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	key, err := sshkey.NewSshKey()
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
		key:         key,
		jupyterPort: *jupyterPort,
		sshPort:     *sshPort,
		containerID: resp.ID,
	}

	err = waitForInstanceReady(ctx, *instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func waitForInstanceReady(ctx context.Context, inst Instance) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	timeout := time.After(1 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return errors.New("jupyter notebook not running")
		case <-ticker.C:
			resp, err := client.Get(fmt.Sprintf("http://localhost:%d", inst.jupyterPort))
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}
