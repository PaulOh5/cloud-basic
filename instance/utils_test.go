package instance

import (
	"fmt"
	"net"
	"testing"
)

func TestAllocatePort(t *testing.T) {
	ports := []*int{new(int), new(int), new(int), new(int)}
	if err := allocatePorts(ports...); err != nil {
		t.Fatal(err)
	}

	portCheckFn := func(port int) {
		addr := net.JoinHostPort("localhost", fmt.Sprintf("%d", port))
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("Port %d is not available: %v", port, err)
		}

		defer listener.Close()
	}

	for _, port := range ports {
		portCheckFn(*port)
	}
}
