package daemon

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/protocol"
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
	files, err := os.ReadDir(d.paths.Services)
	if err != nil {
		return nil
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(d.paths.Services, f.Name())
		cfg, err := service.LoadConfig(path)
		if err != nil {
			// error log here
			continue
		}
		d.registry.Add(&service.State{Config: cfg})
	}

	return nil
}

func (d *Daemon) Dispatch(req *protocol.Request) *protocol.Response {
	switch req.Cmd {
	case "list":
		return d.handleList()
	default:
		return &protocol.Response{
			Ok:    false,
			Error: "unknown command",
		}
	}
}

func (d *Daemon) handleList() *protocol.Response {
	services := d.registry.List()

	data := make(map[string]any)
	svcList := make([]map[string]any, 0, len(services))
	for _, s := range services {
		svcList = append(svcList, map[string]any{
			"name":        s.Config.Name,
			"description": s.Config.Description,
			"author":      s.Config.Author,
			"running":     s.Running,
			"pid":         s.PID,
		})
	}
	data["services"] = svcList

	response := &protocol.Response{
		Ok:   true,
		Data: data,
	}

	log.Printf("responding: %+v", response)

	return response
}
