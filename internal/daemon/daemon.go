package daemon

import (
	"log"
	"net"
	"os"
	"sync"

	"github.com/mhs003/harbrix/internal/helpers"
	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/service"
)

type UserContext struct {
	User     *helpers.User
	Paths    *paths.Paths
	Registry *service.Registry
}

type Daemon struct {
	paths    *paths.Paths
	listener net.Listener

	mu       sync.Mutex
	shutdown bool
	registry *service.Registry
	users    map[string]*UserContext
}

func New() (*Daemon, error) {
	if err := paths.EnsureInternalDir(); err != nil {
		return nil, err
	}

	if err := os.RemoveAll(paths.SocketPath); err != nil {
		return nil, err
	}

	l, err := net.Listen("unix", paths.SocketPath)
	if err != nil {
		return nil, err
	}

	if err := os.Chmod(paths.SocketPath, 0o666); err != nil {
		l.Close()
		return nil, err
	}

	return &Daemon{
		listener: l,
		users:    make(map[string]*UserContext),
	}, nil
}

func (d *Daemon) LoadAllUsers() error {
	opts := helpers.DefaultLoginUserOptions()
	usrs, err := helpers.GetLoginUsers(opts)
	if err != nil {
		return err
	}

	for _, u := range usrs {
		p := paths.NewForHome(u.Home)
		if err := p.Ensure(u.UID, u.GID); err != nil {
			return err
		}

		reg := service.NewRegistry()
		configs, err := service.LoadConfigsFromDisc(p, service.ModeStart)
		if err != nil {
			return err
		}
		for _, cfg := range configs {
			reg.Add(&service.State{
				Config: cfg,
				UID:    u.UID,
				GID:    u.GID,
			})
		}
		d.users[u.Name] = &UserContext{
			User:     &u,
			Paths:    p,
			Registry: reg,
		}
	}
	return nil
}

func (d *Daemon) StartAllEnabled() {
	for _, uc := range d.users {
		entries, err := os.ReadDir(uc.Paths.EnabledServices)
		if err != nil {
			log.Printf("failed reading enabled services for %s: %v", uc.User.Name, err)
			continue
		}
		for _, e := range entries {
			name := e.Name()
			s := uc.Registry.Get(name)
			if s == nil {
				continue
			}
			s.IsEnabled = true
			if err := s.Start(uc.Paths); err != nil {
				log.Printf("failed starting %s/%s: %s", uc.User.Name, name, err)
			}
		}
	}
}

// not used anywhere though.!
func (d *Daemon) ReloadAllUsers() error {
	for _, uc := range d.users {
		configs, err := service.LoadConfigsFromDisc(uc.Paths, service.ModeCLI)
		if err != nil {
			return err
		}
		uc.Registry.Reload(configs)
	}
	return nil
}

func (d *Daemon) ReloadUser(uc *UserContext) error {
	configs, err := service.LoadConfigsFromDisc(uc.Paths, service.ModeCLI)
	if err != nil {
		return err
	}
	uc.Registry.Reload(configs)
	return nil
}
