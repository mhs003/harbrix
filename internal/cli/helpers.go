package cli

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

var bold = "\033[1m"
var green = "\033[32m"
var red = "\033[31m"
var yellow = "\033[33m"
var blue = "\033[34m"
var reset = "\033[0m"

func send(cmd, service string) *protocol.Response {
	p, err := paths.New()
	if err != nil {
		Fatal(err.Error())
	}

	conn, err := net.Dial("unix", p.Socket)
	if err != nil {
		Fatal("daemon not running")
	}
	defer conn.Close()

	req := &protocol.Request{
		Cmd:     cmd,
		Service: service,
	}
	if err := json.NewEncoder(conn).Encode(req); err != nil {
		Fatal(err.Error())
	}

	var resp protocol.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		Fatal(err.Error())
	}

	if !resp.Ok {
		if resp.Error != "" {
			Fatal(resp.Error)
		}
	}

	return &resp
}

func requireArg(args []string) {
	if len(args) < 1 {
		Fatal("service name required")
	}
}

func showLog(name string, follow bool) {
	p, err := paths.New()
	if err != nil {
		Fatal(err.Error())
	}

	logPath := filepath.Join(p.ServiceLogs, name+".log")

	f, err := os.Open(logPath)
	if err != nil {
		Fatal("log file not found")
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
		Fatal(err.Error())
	}

	path := filepath.Join(p.Services, name+".toml")

	if _, err := os.Stat(path); err == nil {
		Fatal("service already exists.")
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
		Fatal(err.Error())
	}

	Success(fmt.Sprintf("created: %s", path))
}

func editService(name string) {
	p, err := paths.New()
	if err != nil {
		Fatal(err.Error())
	}

	path := filepath.Join(p.Services, name+".toml")

	if _, err := os.Stat(path); err != nil {
		Fatal("service not found")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // falling back to `vi`
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Fatal(err.Error())
	}
}

func Fatal(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s%s%s\n", bold, red, msg, reset)
	os.Exit(1)
}

func Success(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s%s\n", bold, green, msg, reset)
	os.Exit(0)
}

func Info(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s%s\n", bold, blue, msg, reset)
	os.Exit(0)
}

func Warning(msg string) {
	fmt.Fprintf(os.Stdout, "%s%s%s%s\n", bold, yellow, msg, reset)
	os.Exit(0)
}
