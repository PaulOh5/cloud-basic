package sshkey

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
)

var KeyDir = os.Getenv("HOME") + "/.ssh/cloud-basic"

func createPrivateKey(keyName string) error {

	// key가 저장되는 디렉토리 생성
	if err := os.MkdirAll(KeyDir, 0755); err != nil {
		return err
	}

	// 빈 파일 생성
	privateKeyPath := KeyDir + "/" + keyName
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	// rsa key 생성
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// rsa key를 pem 형식으로 인코딩
	privatePem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// pem 형식의 rsa key를 파일에 쓰기
	if err := pem.Encode(privateFile, privatePem); err != nil {
		return err
	}

	return nil
}

func removePrivateKey(keyName string) error {
	keyPath := KeyDir + "/" + keyName
	if err := os.Remove(keyPath); err != nil {
		return err
	}
	return nil
}

func createPublicKey(keyName string) error {
	privateKeyPath := KeyDir + "/" + keyName
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	if privateKeyBlock == nil {
		return errors.New("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	return os.WriteFile(privateKeyPath+".pub", ssh.MarshalAuthorizedKey(publicKey), 0644)
}

func removePublicKey(keyName string) error {
	keyPath := KeyDir + "/" + keyName + ".pub"
	if err := os.Remove(keyPath); err != nil {
		return err
	}
	return nil
}

func applyKeyToContainer(cli *client.Client, containerID, keyName string) error {
	publicKeyPath := KeyDir + "/" + keyName + ".pub"
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	publicKey := string(publicKeyFile)

	execFunc := func(cmd []string) error {
		_, err := cli.ContainerExecCreate(
			context.Background(),
			containerID,
			types.ExecConfig{
				Cmd: cmd,
			},
		)
		if err != nil {
			return err
		}
		return nil
	}

	if err := execFunc([]string{"mkdir", "-p", "/root/.ssh"}); err != nil {
		return err
	}

	if err := execFunc([]string{"sh", "-c", "echo " + publicKey + " >> /root/.ssh/authorized_keys"}); err != nil {
		return err
	}

	if err := execFunc([]string{"chmod", "600", "/root/.ssh/authorized_keys"}); err != nil {
		return err
	}

	return nil
}
