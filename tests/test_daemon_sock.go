package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/protocol"
)

func main() {
	p, _ := paths.New()
	conn, _ := net.Dial("unix", p.Socket)
	defer conn.Close()

	req := &protocol.Request{
		Cmd: "list",
	}
	json.NewEncoder(conn).Encode(req)

	var resp protocol.Response
	json.NewDecoder(conn).Decode(&resp)
	fmt.Printf("response: %+v\n", resp)
}
