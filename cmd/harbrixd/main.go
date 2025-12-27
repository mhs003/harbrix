package main

import (
	"log"
	"os"

	"github.com/mhs003/harbrix/internal/daemon"
	"github.com/mhs003/harbrix/internal/paths"
)

func main() {
	p, err := paths.New()
	if err != nil {
		log.Fatalf("paths init failed: %v", err)
	}

	if err := p.Ensure(); err != nil {
		log.Fatalf("failed to ensure directories: %v", err)
	}

	// init and run the daemon
	d, err := daemon.New(p)
	if err != nil {
		log.Fatalf("daemon init failed: %v", err)
	}

	if err := d.LoadServices(); err != nil {
		log.Fatalf("failed loading services: %v", err)
	}

	if err := d.Run(); err != nil {
		log.Printf("daemon stopped: %v", err)
		os.Exit(1)
	}
}
