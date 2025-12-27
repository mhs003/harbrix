package daemon

import "net"

func (d *Daemon) Run() error {
	for {
		conn, err := d.listener.Accept()
		if err != nil {
			if d.isShuttingDown() {
				return nil
			}
			continue
		}

		go d.handleConn(conn)
	}
}

func (d *Daemon) handleConn(conn net.Conn) {
	defer conn.Close()
	// TODO: handle ipc connection here
}

func (d *Daemon) isShuttingDown() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.shutdown
}
