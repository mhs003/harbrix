package daemon

import "github.com/mhs003/harbrix/internal/protocol"

// the dispatcher
func (d *Daemon) Dispatch(req *protocol.Request) *protocol.Response {
	switch req.Cmd {
	case "start":
		return d.handleStart(req.Service)
	case "stop":
		return d.handleStop(req.Service)
	case "list":
		return d.handleList()
	case "reload-daemon":
		return d.handleReload()
	default:
		return &protocol.Response{
			Ok:    false,
			Error: "unknown command",
		}
	}
}
