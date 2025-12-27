package service

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Author      string `toml:"author"`

	Service struct {
		Command string `toml:"command"`
		Workdir string `toml:"workdir"`
		Restart string `toml:"restart"`
	} `toml:"service"`
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("service file not found.")
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	if filepath.Base(path) != cfg.Name+".toml" {
		return nil, errors.New("filename must match service name")
	}

	// basic validation
	if cfg.Service.Command == "" {
		return nil, errors.New("service.command cannot be empty")
	}

	switch cfg.Service.Restart {
	case "", "never", "on-failure", "always":
	default:
		return nil, errors.New("invalid restart policy")
	}

	return &cfg, nil
}
