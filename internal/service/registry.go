package service

import (
	"errors"
	"os/exec"
	"sync"
)

type State struct {
	Config    *Config
	Running   bool
	PID       int
	Cmd       *exec.Cmd
	ExitCode  int
	StopReq   bool // manual stop requested ; used to prevent auto restart on user stop
	IsEnabled bool
}

type Registry struct {
	mu       sync.Mutex
	services map[string]*State
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*State),
	}
}

func (r *Registry) Add(s *State) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[s.Config.Name]; exists {
		return errors.New("service already exists in registry")
	}
	r.services[s.Config.Name] = s
	return nil
}

func (r *Registry) List() []*State {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := make([]*State, 0, len(r.services))
	for _, s := range r.services {
		list = append(list, s)
	}

	return list
}

func (r *Registry) Get(name string) *State {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.services[name]
}

// func (r *Registry) Update(name string, )

func (r *Registry) Reload(new map[string]*Config) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, cfg := range new {
		s, exists := r.services[name]
		if !exists {
			r.services[name] = &State{
				Config:  cfg,
				Running: false,
				PID:     0,
			}
			continue
		}

		if s.Running {
			continue
		}

		s.Config = cfg
	}

	for name, s := range r.services {
		if _, exists := new[name]; !exists && !s.Running {
			delete(r.services, name)
		}
	}
}
