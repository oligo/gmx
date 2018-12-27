package gmx

import (
	"net"
	"fmt"
	"path/filepath"
	"os"
	"strconv"
	"runtime"
)

// SockType defines gmx connection socket type
type SockType int

const (
	SOCK_UNIX  SockType = 0
	SOCK_TCP SockType = 1
)

const (
	//EnvGMXHost is environment variable name for tcp socket host
	EnvGMXHost = "GMX_HOST"
	//EnvGMXPort is environment variable name for tcp socket port
	EnvGMXPort = "GMX_PORT"
)

const (
	DEFAULT_BIND_IP = "127.0.0.1"
	DEFAULT_PORT = 9997
)

func tcpSocket(addr string, port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
}

func localSocket(args ...interface{}) (net.Listener, error) {
	return net.ListenUnix("unix", localSocketAddr())
}

func localSocketAddr() *net.UnixAddr {
	return &net.UnixAddr{
		filepath.Join(os.TempDir(), fmt.Sprintf(".gmx.%d.%d", os.Getpid(), GMX_VERSION)),
		"unix",
	}
}

func setupSocket() (net.Listener, error) {
	var sockChoice = SOCK_UNIX
	var host string
	var port int

	if runtime.GOOS == "windows" {
		sockChoice = SOCK_TCP
	} 

	host = os.Getenv(EnvGMXHost)
	if (len(host) > 0) {
		sockChoice = SOCK_TCP
		var err error
		if port, err = strconv.Atoi(os.Getenv(EnvGMXPort)); err != nil {
			// parse error or not provided
			port = DEFAULT_PORT
		}
	} else if sockChoice == SOCK_TCP {
		host = DEFAULT_BIND_IP
		port = DEFAULT_PORT
	}

	
	switch sockChoice {
	case SOCK_TCP:
		return tcpSocket(host, port)
	default:
		return localSocket()
	}
}  