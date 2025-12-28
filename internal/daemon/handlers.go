package daemon

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mhs003/harbrix/internal/protocol"
)

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
			"enabled":     s.IsEnabled,
			// "cmd":         fmt.Sprintf("%+v", s.Cmd.),
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

func (d *Daemon) handleStart(name string) *protocol.Response {
	s := d.registry.Get(name)
	if s == nil {
		return &protocol.Response{Ok: false, Error: "service not found"}
	}

	if err := s.Start(d.paths); err != nil {
		return &protocol.Response{Ok: false, Error: err.Error()}
	}

	return &protocol.Response{Ok: true}
}

func (d *Daemon) handleStop(name string) *protocol.Response {
	s := d.registry.Get(name)
	if s == nil {
		return &protocol.Response{Ok: false, Error: "service not found"}
	}

	if err := s.Stop(); err != nil {
		return &protocol.Response{Ok: false, Error: err.Error()}
	}

	return &protocol.Response{Ok: true}
}

// func (d *Daemon) handleRestart(name string) *protocol.Response {
// 	s := d.registry.Get(name)
// 	if s == nil {
// 		return &protocol.Response{Ok: false, Error: "service not found"}
// 	}
// 	if err := s.Stop(); err != nil {
// 		return &protocol.Response{Ok: false, Error: err.Error()}
// 	}
// 	if err := s.Start(d.paths); err != nil {
// 		return &protocol.Response{Ok: false, Error: err.Error()}
// 	}
// 	return &protocol.Response{Ok: true}
// }

func (d *Daemon) handleReload() *protocol.Response {
	if err := d.ReloadServices(); err != nil {
		return &protocol.Response{Ok: false, Error: err.Error()}
	}
	return &protocol.Response{Ok: true}
}

func (d *Daemon) handleEnable(name string) *protocol.Response {
	s := d.registry.Get(name)
	if s == nil {
		return &protocol.Response{Ok: false, Error: "service not found"}
	}

	path := filepath.Join(d.paths.EnabledServices, name)
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		return &protocol.Response{Ok: false, Error: err.Error()}
	}

	s.IsEnabled = true

	return &protocol.Response{Ok: true}
}

func (d *Daemon) handleDisable(name string) *protocol.Response {
	path := filepath.Join(d.paths.EnabledServices, name)
	os.Remove(path)
	if s := d.registry.Get(name); s != nil {
		s.IsEnabled = false
	}
	return &protocol.Response{Ok: true}
}

func (d *Daemon) handleIsEnabled(name string) *protocol.Response {
	path := filepath.Join(d.paths.EnabledServices, name)
	_, err := os.Stat(path)
	if err != nil {
		return &protocol.Response{Ok: false}
	}
	return &protocol.Response{Ok: true}
}
