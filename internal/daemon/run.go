package daemon

import (
	"log"
	"net"

	"github.com/mhs003/harbrix/internal/protocol"
)

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
	req, err := protocol.DecodeRequest(conn)
	if err != nil {
		log.Printf("decode error: %v", err)
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "invalid request",
		})
		return
	}

	log.Printf("received request: cmd=%s service=%s", req.Cmd, req.Service)

	resp := d.Dispatch(req)
	if err := protocol.EncodeResponse(conn, resp); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (d *Daemon) isShuttingDown() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.shutdown
}
