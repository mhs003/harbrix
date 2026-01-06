package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mhs003/harbrix/internal/paths"
)

type Config struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Author      string `toml:"author"`

	Env map[string]any `toml:"env"`

	Service struct {
		Command string `toml:"command"`
		Workdir string `toml:"workdir"`
		Log     bool   `toml:"log"`
	} `toml:"service"`

	Restart struct {
		Policy    string `toml:"policy"`
		Delay     string `toml:"delay"`     // default=0s
		Limit     int    `toml:"limit"`     // no limit by default
		MaxFailed int    `toml:"maxfailed"` // default will be 5, set to -1 to bypass limitation
	} `toml:"restart"`
}

func ApplyDefaults(c *Config) {
	if c.Restart.Policy == "" {
		c.Restart.Policy = "never"
	}
	if c.Restart.Delay == "" {
		c.Restart.Delay = "0s"
	}
	if c.Restart.MaxFailed == 0 {
		c.Restart.MaxFailed = 5
	}
}

func (c *Config) ValidateConfig() error {
	if c.Service.Command == "" {
		return errors.New("service command cannot be empty")
	}

	switch c.Restart.Policy {
	case "never", "on-failure", "always":
	default:
		return errors.New("invalid restart policy")
	}

	for k := range c.Env {
		if k == "" {
			return errors.New("env key cannot be empty")
		}
	}

	if _, err := time.ParseDuration(c.Restart.Delay); err != nil {
		return errors.New("invalid restart.delay")
	}

	if c.Restart.Limit < 0 {
		return errors.New("invalid restart.limit")
	}

	if c.Restart.MaxFailed < -1 {
		return errors.New("invalid restart.maxfailed")
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("no services found at path: %s", path)
		return nil, fmt.Errorf("no services found at path: %s", path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		log.Printf("caught error loading service config: %s", err)
		return nil, err
	}

	if filepath.Base(path) != cfg.Name+".toml" {
		log.Printf("filename must match service name for service %s", path)
		return nil, fmt.Errorf("filename must match service name for service %s", path)
	}

	ApplyDefaults(&cfg)

	return &cfg, nil
}

func LoadConfigsFromDisc(paths *paths.Paths, mode LoadMode) (map[string]*Config, error) {
	cfgs := make(map[string]*Config)

	files, err := os.ReadDir(paths.Services)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(paths.Services, f.Name())
		cfg, err := LoadConfig(path)
		if err != nil {
			if mode == ModeCLI {
				return nil, fmt.Errorf("failed to load %s: %w", path, err)
			} else {
				log.Printf("invalid config %s: %v", path, err)
				continue
			}
		}

		cfgs[cfg.Name] = cfg
	}

	return cfgs, nil
}
