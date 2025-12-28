package daemon

import "github.com/mhs003/harbrix/internal/protocol"

func (d *Daemon) Dispatch(req *protocol.Request) *protocol.Response {
	switch req.Cmd {
	case "start":
		return d.handleStart(req.Service)
	case "stop":
		return d.handleStop(req.Service)
	// case "restart":
	// 	return d.handleRestart(req.Service)
	case "list":
		return d.handleList()
	case "reload-daemon":
		return d.handleReload()
	case "enable":
		return d.handleEnable(req.Service)
	case "disable":
		return d.handleDisable(req.Service)
	case "is-enabled":
		return d.handleIsEnabled(req.Service)
	default:
		return &protocol.Response{
			Ok:    false,
			Error: "unknown command",
		}
	}
}
