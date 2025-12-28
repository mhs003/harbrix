package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

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
	case "log":
		if arg == "" {
			fatal("service name required")
		}

		follow := false
		if len(os.Args) > 2 && os.Args[2] == "-f" {
			if len(os.Args) < 4 {
				fatal("service name required")
			}
			follow = true
			arg = os.Args[3]
		}
		showLog(arg, follow)
	case "reload-daemon":
		send("reload-daemon", "")
	default:
		fatal(fmt.Sprintf("unknown command %+v", os.Args))
	}
}

func showLog(name string, follow bool) {
	p, err := paths.New()
	if err != nil {
		fatal(err.Error())
	}

	logPath := filepath.Join(p.ServiceLogs, name+".log")

	f, err := os.Open(logPath)
	if err != nil {
		fatal("log file not found")
	}
	defer f.Close()

	if !follow {
		io.Copy(os.Stdout, f)
		return
	}

	f.Seek(0, io.SeekEnd)
	buf := make([]byte, 4096)

	for {
		n, err := f.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
		}
		if err != nil {
			time.Sleep(300 * time.Millisecond)
		}
	}
}
