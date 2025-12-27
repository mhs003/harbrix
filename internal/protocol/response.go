package protocol

import (
	"encoding/json"
	"io"
)

type Response struct {
	Ok    bool           `json:"ok"`
	Error string         `json:"error,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

func EncodeResponse(w io.Writer, resp *Response) error {
	return json.NewEncoder(w).Encode(resp)
}
