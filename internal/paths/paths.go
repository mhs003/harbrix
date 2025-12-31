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
	// Socket          string
}

const InternalDir = "/run/harbrix"
const SocketPath = InternalDir + "/control.sock"

func New() (*Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return NewForHome(home), nil
}

func NewForHome(home string) *Paths {
	root := filepath.Join(home, ".local", "share", "harbrix")

	return &Paths{
		Root:            root,
		Services:        filepath.Join(root, "services"),
		Logs:            filepath.Join(root, "logs"),
		ServiceLogs:     filepath.Join(root, "logs", "services"),
		State:           filepath.Join(root, "state"),
		EnabledServices: filepath.Join(root, "enabled"),
	}
}

func EnsureInternalDir() error {
	if err := os.MkdirAll(InternalDir, 0o755); err != nil {
		return err
	}
	return nil
}

func (p *Paths) Ensure(uid int, gid int) error {
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
		if err := os.Chown(dir, uid, gid); err != nil {
			return err
		}
	}

	return nil
}
