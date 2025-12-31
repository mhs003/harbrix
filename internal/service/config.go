package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mhs003/harbrix/internal/paths"
)

type Config struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Author      string `toml:"author"`

	Service struct {
		Command string `toml:"command"`
		Workdir string `toml:"workdir"`
		Restart string `toml:"restart"`
		Log     bool   `toml:"log"`
	} `toml:"service"`
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

	return &cfg, nil
}

func LoadConfigsFromDisc(paths *paths.Paths) (map[string]*Config, error) {
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
		cfg, _ := LoadConfig(path)
		// if err != nil {
		// 	return nil, err
		// }
		cfgs[cfg.Name] = cfg
	}

	return cfgs, nil
}
