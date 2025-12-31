package daemon

import (
	"log"
	"net"
	"syscall"

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

	uconn, ok := conn.(*net.UnixConn)
	if !ok {
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "invalid connection",
		})
		return
	}
	f, err := uconn.File()
	if err != nil {
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "failed to get connection file",
		})
		return
	}
	defer f.Close()

	ucred, err := syscall.GetsockoptUcred(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	if err != nil {
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "failed to get peer credential",
		})
		return
	}

	callerID := int(ucred.Uid)

	var uc *UserContext
	for _, uctx := range d.users {
		if uctx.User.UID == callerID {
			uc = uctx
			break
		}
	}
	if uc == nil {
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "unauthorized user",
		})
		return
	}

	if req.Env == nil || req.Env["USER"] != uc.User.Name || req.Env["HOME"] != uc.User.Home {
		protocol.EncodeResponse(conn, &protocol.Response{
			Ok:    false,
			Error: "invalid environment",
		})
		return
	}

	d.paths = uc.Paths
	d.registry = uc.Registry

	log.Printf("received request from %s: cmd=%s service=%s", uc.User.Name, req.Cmd, req.Service)

	resp := d.Dispatch(req, uc)
	if err := protocol.EncodeResponse(conn, resp); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func (d *Daemon) isShuttingDown() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.shutdown
}
