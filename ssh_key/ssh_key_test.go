package sshkey

import (
	"os"
	"testing"
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

}
