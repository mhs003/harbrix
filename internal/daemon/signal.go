package daemon

import (
	"os"
	"os/signal"
	"syscall"
)

func (d *Daemon) InitSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		d.Shutdown()
	}()
}

func (d *Daemon) Shutdown() {
	d.mu.Lock()
	d.shutdown = true
	d.mu.Unlock()

	d.listener.Close()
	_ = os.Remove(d.paths.Socket)
}
