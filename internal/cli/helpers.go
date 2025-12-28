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
		Fatal(resp.Error)
	}

	return &resp
}

func Fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
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

	fmt.Println("created:", path)
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
