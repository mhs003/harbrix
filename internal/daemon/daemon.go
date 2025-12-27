package daemon

import (
	"net"
	"os"
	"sync"

	"github.com/mhs003/harbrix/internal/paths"
)

type Daemon struct {
	paths    *paths.Paths
	listener net.Listener

	mu       sync.Mutex
	shutdown bool
}

func New(p *paths.Paths) (*Daemon, error) {
	if err := os.RemoveAll(p.Socket); err != nil {
		return nil, err
	}

	l, err := net.Listen("unix", p.Socket)
	if err != nil {
		return nil, err
	}

	if err := os.Chmod(p.Socket, 0o600); err != nil {
		l.Close()
		return nil, err
	}

	return &Daemon{
		paths:    p,
		listener: l,
	}, nil
}
