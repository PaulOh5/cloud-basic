package instance

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PaulOh5/cloud-basic/docker"
	"github.com/PaulOh5/cloud-basic/network"
	sshkey "github.com/PaulOh5/cloud-basic/ssh_key"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/crypto/ssh"
)

type Instance struct {
	cli                  *client.Client
	key                  *sshkey.SSHKey
	jupyterPort, sshPort int
	jupyterURL           string
	containerID          string
}

func (inst *Instance) Exec(ctx context.Context, cmd ...string) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execID, err := inst.cli.ContainerExecCreate(ctx, inst.containerID, execConfig)
	if err != nil {
		return "", err
	}

	resp, err := inst.cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	var outBuf, errBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
	if err != nil {
		return "", err
	}

	if errBuf.Len() > 0 {
		return outBuf.String(), errors.New(errBuf.String())
	}

	return outBuf.String(), nil
}

func (inst *Instance) Start(ctx context.Context) error {
	err := inst.cli.ContainerStart(ctx, inst.containerID, container.StartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (inst *Instance) Stop(ctx context.Context) error {
	err := inst.cli.ContainerStop(ctx, inst.containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (inst *Instance) Remove(ctx context.Context) error {
	err := inst.cli.ContainerRemove(context.Background(), inst.containerID, container.RemoveOptions{Force: true})
	if err != nil {
		return err
	}

	inst.cli.Close()
	return nil
}

func (inst *Instance) GetStatus(ctx context.Context) (InstanceStatus, error) {
	ctr, err := inst.cli.ContainerInspect(ctx, inst.containerID)
	if err != nil {
		return STOPPED, err
	}

	if ctr.State.Running {
		return RUNNING, nil
	}

	return STOPPED, nil
}

func (inst *Instance) EstablishConnect(fh *network.ForwardingHandler) error {
	jupyterPath := fmt.Sprintf("http://127.0.0.1:%d", inst.jupyterPort)
	err := fh.AddForwarding("/"+inst.containerID[12:]+"/jupyter", jupyterPath)
	if err != nil {
		return err
	}

	inst.jupyterURL = "/" + inst.containerID[12:] + "/jupyter"
	return nil
}

func (inst *Instance) Disconnect(fh *network.ForwardingHandler) {
	fh.RemoveForwarding(inst.jupyterURL)
}

func (inst *Instance) GetSshUrl() string {
	return fmt.Sprintf("localhost:%d", inst.sshPort)
}

func (inst *Instance) GetSshConfig() (*ssh.ClientConfig, error) {
	if inst.key == nil {
		return nil, errors.New("ssh key not set")
	}

	signer, err := ssh.NewSignerFromKey(inst.key.PrivateKey)
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

func (inst *Instance) SetSshKey() error {
	key, err := sshkey.NewSshKey()
	if err != nil {
		return err
	}

	cmdGroup := [][]string{
		{"mkdir", "-p", "/root/.ssh"},
		{"echo", "\"" + string(ssh.MarshalAuthorizedKey(key.PublicKey)) + "\"", ">", "/root/.ssh/authorized_keys"},
		{"chmod", "600", "/root/.ssh/authorized_keys"},
	}

	fmt.Println(cmdGroup[1])

	for _, cmd := range cmdGroup {
		_, err = inst.Exec(context.Background(), cmd...)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: 리팩토링 필요
func NewInstance(ctx context.Context) (*Instance, error) {
	instId, jPort, sPort, cli, err := docker.GenerateCloudContainer(ctx)
	if err != nil {
		return nil, err
	}

	inst := &Instance{
		cli:         cli,
		jupyterPort: jPort,
		sshPort:     sPort,
		containerID: instId,
	}

	err = inst.SetSshKey()
	if err != nil {
		inst.Remove(ctx)
		return nil, err
	}

	err = waitForInstanceReady(ctx, inst)
	if err != nil {
		inst.Remove(ctx)
		return nil, err
	}

	return inst, nil
}

func waitForInstanceReady(ctx context.Context, inst *Instance) error {
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
