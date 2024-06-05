package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"os"

	"golang.org/x/crypto/ssh"
)

var KeyDir = os.Getenv("HOME") + "/.ssh/cloud-basic"

type SSHKey struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  ssh.PublicKey
}

func NewSshKey() (*SSHKey, error) {
	key := &SSHKey{}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	key.PrivateKey = privateKey
	key.PublicKey = publicKey

	return key, nil
}

// func applyKeyToContainer(cli *client.Client, containerID, keyName string) error {
// 	publicKeyPath := KeyDir + "/" + keyName + ".pub"
// 	publicKeyFile, err := os.ReadFile(publicKeyPath)
// 	if err != nil {
// 		return err
// 	}

// 	publicKey := string(publicKeyFile)

// 	execFunc := func(cmd []string) error {
// 		_, err := cli.ContainerExecCreate(
// 			context.Background(),
// 			containerID,
// 			types.ExecConfig{
// 				Cmd: cmd,
// 			},
// 		)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	if err := execFunc([]string{"mkdir", "-p", "/root/.ssh"}); err != nil {
// 		return err
// 	}

// 	if err := execFunc([]string{"sh", "-c", "echo " + publicKey + " >> /root/.ssh/authorized_keys"}); err != nil {
// 		return err
// 	}

// 	if err := execFunc([]string{"chmod", "600", "/root/.ssh/authorized_keys"}); err != nil {
// 		return err
// 	}

// 	return nil
// }
