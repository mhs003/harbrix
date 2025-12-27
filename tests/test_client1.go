package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/mhs003/harbrix/internal/paths"
	"github.com/mhs003/harbrix/internal/protocol"
)

func sendRequest(p *paths.Paths, req *protocol.Request) *protocol.Response {
	conn, err := net.Dial("unix", p.Socket)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		panic(err)
	}
	var resp protocol.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		panic(err)
	}

	return &resp
}

func main() {
	p, _ := paths.New()

	fmt.Println("=== LIST services ===")
	resp := sendRequest(p, &protocol.Request{Cmd: "list"})
	fmt.Printf("%+v\n", resp)

	fmt.Println("=== START 'test' ===")
	resp = sendRequest(p, &protocol.Request{Cmd: "start", Service: "test"})
	fmt.Printf("%+v\n", resp)

	time.Sleep(1 * time.Second)

	fmt.Println("=== LIST services ===")
	resp = sendRequest(p, &protocol.Request{Cmd: "list"})
	fmt.Printf("%+v\n", resp)

	fmt.Println("=== STOP 'test' ===")
	resp = sendRequest(p, &protocol.Request{Cmd: "stop", Service: "test"})
	fmt.Printf("%+v\n", resp)

	fmt.Println("=== LIST services ===")
	resp = sendRequest(p, &protocol.Request{Cmd: "list"})
	fmt.Printf("%+v\n", resp)
}
