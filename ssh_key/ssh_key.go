package sshkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

var SSHKeyDir = os.Getenv("HOME") + "/.ssh/cloud-basic"

func createSSHPrivateKey(keyName string) error {
	if err := os.MkdirAll(SSHKeyDir, 0755); err != nil {
		return err
	}

	privateKeyPath := SSHKeyDir + "/" + keyName
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

func removeSSHKey(keyName string) error {
	keyPath := SSHKeyDir + "/" + keyName
	if err := os.Remove(keyPath); err != nil {
		return err
	}
	return nil
}
