package main

import (
	"log"
	"os"

	"github.com/mhs003/harbrix/internal/daemon"
)

func main() {
	d, err := daemon.New()
	if err != nil {
		log.Fatalf("daemon init failed: %v", err)
	}
	d.InitSignals()
	if err := d.LoadAllUsers(); err != nil {
		log.Fatalf("load user failed: %v", err)
	}
	d.StartAllEnabled()
	if err := d.Run(); err != nil {
		log.Printf("daemon stopped: %v", err)
		os.Exit(1)
	}
}
