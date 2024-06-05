package sshkey

import (
	"testing"
)

func TestSshKey(t *testing.T) {
	key, err := NewSshKey()
	if err != nil {
		t.Fatal(err)
	}

	if key.PrivateKey == nil {
		t.Fatal("private key not created")
	}

	if key.PublicKey == nil {
		t.Fatal("public key not created")
	}
}
