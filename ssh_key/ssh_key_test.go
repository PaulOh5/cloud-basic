package sshkey

import (
	"os"
	"testing"

	"github.com/PaulOh5/cloud-basic/instance"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
)

func TestCreateKey(t *testing.T) {
	testCases := []struct {
		keyName string
	}{
		{keyName: "paul"},
	}

	for _, tc := range testCases {
		t.Run(tc.keyName, func(t *testing.T) {
			err := createPrivateKey(tc.keyName)
			if err != nil {
				t.Fatal(err)
			}
			defer removePrivateKey(tc.keyName)

			if _, err := os.Stat(KeyDir + "/" + tc.keyName); os.IsNotExist(err) {
				t.Fatalf("ssh private key not created")
			}

			err = createPublicKey(tc.keyName)
			if err != nil {
				t.Fatal(err)
			}
			defer removePublicKey(tc.keyName)

			if _, err := os.Stat(KeyDir + "/" + tc.keyName + ".pub"); os.IsNotExist(err) {
				t.Fatalf("ssh public key not created")
			}
		})
	}
}

func TestApplyKeyToContainer(t *testing.T) {
	testCases := []struct {
		keyName string
	}{
		{keyName: "paul"},
	}

	for _, tc := range testCases {
		t.Run(tc.keyName, func(t *testing.T) {
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				t.Fatal(err)
			}
			defer cli.Close()

			config := instance.ContainerConfig{
				JupyterPort: "9999",
				SshPort:     "2222",
			}
			containerID, err := instance.CreateContainer(cli, config)
			if err != nil {
				t.Fatal(err)
			}
			defer instance.RemoveContainer(cli, containerID)

			err = createPrivateKey(tc.keyName)
			if err != nil {
				t.Fatal(err)
			}
			defer removePrivateKey(tc.keyName)

			err = createPublicKey(tc.keyName)
			if err != nil {
				t.Fatal(err)
			}
			defer removePublicKey(tc.keyName)

			err = applyKeyToContainer(cli, containerID, tc.keyName)
			if err != nil {
				t.Fatal(err)
			}

			if err := checkKeyAppliedToContainer(tc.keyName, config); err != nil {
				t.Fatalf("failed to apply key to container: %v", err)
			}
		})
	}
}

func checkKeyAppliedToContainer(keyName string, containerConfig instance.ContainerConfig) error {
	keyPath := KeyDir + "/" + keyName
	keyFile, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(keyFile)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", "localhost:"+containerConfig.SshPort, sshConfig)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return nil
}
