package daemon

import (
	"log"

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
			// "cmd":         fmt.Sprintf("%+v", s.Cmd),
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
