package service

import (
	"errors"
	"sync"
)

type State struct {
	Config  *Config
	Running bool
	PID     int
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
