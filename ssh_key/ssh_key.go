package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"golang.org/x/crypto/ssh"
)

var KeyDir = os.Getenv("HOME") + "/.ssh/cloud-basic"

func createPrivateKey(keyName string) error {
	if err := os.MkdirAll(KeyDir, 0755); err != nil {
		return err
	}

	privateKeyPath := KeyDir + "/" + keyName
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privatePem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

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
