package daemon

import (
	"net"
	"os"
	"sync"

	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/service"
)

type Daemon struct {
	paths    *paths.Paths
	listener net.Listener

	mu       sync.Mutex
	shutdown bool
	registry *service.Registry
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
		registry: service.NewRegistry(),
	}, nil
}

func (d *Daemon) LoadServices() error {
	configs, err := service.LoadConfigsFromDisc(d.paths)
	if err != nil {
		return err
	}

	for _, cfg := range configs {
		d.registry.Add(&service.State{
			Config: cfg,
		})
	}

	return nil
}

func (d *Daemon) ReloadServices() error {
	configs, err := service.LoadConfigsFromDisc(d.paths)
	if err != nil {
		return err
	}

	d.registry.Reload(configs)
	return nil
}

func (d *Daemon) StartEnabled() {
	entries, err := os.ReadDir(d.paths.EnabledServices)
	if err != nil {
		return
	}

	for _, e := range entries {
		name := e.Name()
		s := d.registry.Get(name)
		if s == nil {
			continue
		}
		s.Start(d.paths)
	}
}
