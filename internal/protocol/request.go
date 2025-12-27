package protocol

import (
	"encoding/json"
	"io"
)

type Request struct {
	Cmd     string                 `json:"cmd"`
	Service string                 `json:"service,omitempty"`
	Args    map[string]interface{} `json:"args"`
}

func DecodeRequest(r io.Reader) (*Request, error) {
	var req Request
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
