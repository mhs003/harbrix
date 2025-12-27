package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/protocol"
)

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}

func send(cmd, service string) *protocol.Response {
	p, err := paths.New()
	if err != nil {
		fatal(err.Error())
	}

	conn, err := net.Dial("unix", p.Socket)
	if err != nil {
		fatal("daemon not running")
	}
	defer conn.Close()

	req := &protocol.Request{
		Cmd:     cmd,
		Service: service,
	}
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		fatal(err.Error())
	}

	var resp protocol.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		fatal(err.Error())
	}

	if !resp.Ok {
		fatal(resp.Error)
	}

	return &resp
}

func main() {
	if len(os.Args) < 2 {
		fatal("no command")
	}

	cmd := os.Args[1]
	arg := ""
	if len(os.Args) > 2 {
		arg = os.Args[2]
	}

	switch cmd {
	case "list":
		resp := send("list", "")
		services := resp.Data["services"].([]interface{})

		for _, raw := range services {
			s := raw.(map[string]interface{})
			status := "stopped"
			if s["running"].(bool) {
				status = fmt.Sprintf("running (pid %v)", s["pid"])
			}
			fmt.Printf("%-10s %s\n", s["name"], status)
		}
	case "start", "stop", "restart":
		if arg == "" {
			fatal("service name required")
		}

		if cmd == "restart" {
			send("stop", arg)
			send("start", arg)
		} else {
			send(cmd, arg)
		}
	default:
		fatal(fmt.Sprintf("unknown command %+v", os.Args))
	}
}
