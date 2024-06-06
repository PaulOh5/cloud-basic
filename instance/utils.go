package instance

import (
	"net"
)

func getAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func allocatePorts(ports ...*int) error {
	for _, port := range ports {
		p, err := getAvailablePort()
		if err != nil {
			return err
		}
		*port = p
	}

	return nil
}
