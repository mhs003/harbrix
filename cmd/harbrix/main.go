package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"
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
	case "enable":
		if arg == "" {
			fatal("service name is required")
		}
		send("enable", arg)
		if len(os.Args) > 3 && os.Args[3] == "--now" {
			send("start", arg)
		}
	case "disable", "is-enabled":
		if arg == "" {
			fatal("service name is required")
		}
		resp := send(cmd, arg)

		if cmd == "is-enabled" {
			if resp.Ok {
				os.Exit(0)
			}
			os.Exit(1)
		}
	case "new":
		if arg == "" {
			fatal("service name is required")
		}
		createService(arg)
	case "edit":
		if arg == "" {
			fatal("service name is required")
		}
		editService(arg)
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

func createService(name string) {
	p, err := paths.New()
	if err != nil {
		fatal(err.Error())
	}

	path := filepath.Join(p.Services, name+".toml")

	if _, err := os.Stat(path); err == nil {
		fatal("service already exists.")
	}

	username := ""
	if u, err := user.Current(); err == nil {
		username = u.Username
	}

	template := fmt.Sprintf(`# %s service

name="%s"
description=""
auth="%s"

[service]
command=""
workdir=""
restart="never"
`, name, name, username)

	if err := os.WriteFile(path, []byte(template), 0o644); err != nil {
		fatal(err.Error())
	}

	fmt.Println("created:", path)
}

func editService(name string) {
	p, err := paths.New()
	if err != nil {
		fatal(err.Error())
	}

	path := filepath.Join(p.Services, name+".toml")

	if _, err := os.Stat(path); err != nil {
		fatal("service not found")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fatal(err.Error())
	}
}
