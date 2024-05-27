package sshkey

import (
	"os"
	"testing"
)

func TestCreateSSHKey(t *testing.T) {
	testCases := []struct {
		keyName string
	}{
		{keyName: "paul"},
	}

	for _, tc := range testCases {
		t.Run(tc.keyName, func(t *testing.T) {
			err := createSSHPrivateKey(tc.keyName)
			if err != nil {
				t.Fatal(err)
			}
			defer removeSSHKey(tc.keyName)

			if _, err := os.Stat(SSHKeyDir + "/" + tc.keyName); os.IsNotExist(err) {
				t.Fatalf("ssh private key not created")
			}
		})
	}
}
