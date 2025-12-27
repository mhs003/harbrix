package main

import (
	"log"

	"github.com/mhs003/harbrix/internal/paths"
)

func main() {
	p, err := paths.New()
	if err != nil {
		log.Fatalf("paths init failed: %v\n", err)
	}

	if err := p.Ensure(); err != nil {
		log.Fatalf("failed to ensure directories: %v\n", err)
	}

	// create and run the daemon
}
