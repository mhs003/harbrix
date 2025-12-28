package paths

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Root            string
	Services        string
	Logs            string
	ServiceLogs     string
	State           string
	EnabledServices string
	Socket          string
}

func New() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(home, ".local", "share", "harbrix")

	return &Paths{
		Root:            root,
		Services:        filepath.Join(root, "services"),
		Logs:            filepath.Join(root, "logs"),
		ServiceLogs:     filepath.Join(root, "logs", "services"),
		State:           filepath.Join(root, "state"),
		EnabledServices: filepath.Join(root, "enabled"),
		Socket:          filepath.Join(root, "control.sock"),
	}, nil
}

func (p *Paths) Ensure() error {
	dirs := []string{
		p.Root,
		p.Services,
		p.Logs,
		p.ServiceLogs,
		p.State,
		p.EnabledServices,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}
