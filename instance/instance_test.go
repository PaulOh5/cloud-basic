package instance

import (
	"context"
	"testing"
	"time"
)

func TestInstance(t *testing.T) {
	instance, err := NewInstance(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	running, err := instance.IsRunning(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !running {
		t.Fatal("instance not running")
	}
}

// func TestContainer(t *testing.T) {
// 	testCases := []struct {
// 		config ContainerConfig
// 		key    *sshkey.SSHKey
// 	}{
// 		{config: ContainerConfig{
// 			JupyterPort: "9999",
// 			SshPort:     "2222",
// 		}},
// 	}

// 	for i, tc := range testCases {
// 		t.Run(tc.config.JupyterPort, func(t *testing.T) {
// 			tc.key, _ = sshkey.NewSshKey()

// 			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			defer cli.Close()
// 			containerID, err := CreateContainer(cli, tc.config)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			defer RemoveContainer(cli, containerID)
// 			time.Sleep(5 * time.Second)

// 			if err := checkContainerCreated(cli, containerID); err != nil {
// 				t.Fatalf("container not created: %v", err)
// 			}
// 			if err := checkJupyterResponse(tc.config.JupyterPort); err != nil {
// 				t.Fatalf("failed to connect to jupyter notebook: %v", err)
// 			}
// 			if err := checkSshConnection(tc.config); err != nil {
// 				t.Fatalf("failed to connect to ssh: %v", err)
// 			}
// 			t.Logf("test case %d passed", i+1)
// 		})
// 	}
// }

// func checkContainerCreated(cli *client.Client, containerID string) error {
// 	ctr, err := cli.ContainerInspect(context.Background(), containerID)
// 	if err != nil {
// 		return err
// 	}

// 	if ctr.ID != containerID {
// 		return fmt.Errorf("container not found")
// 	}

// 	if !ctr.State.Running {
// 		return fmt.Errorf("container not running")
// 	}
// 	return nil
// }

// func checkJupyterResponse(port string) error {
// 	resp, err := http.Get("http://localhost:" + port)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("jupyter notebook not running")
// 	}
// 	return nil
// }
