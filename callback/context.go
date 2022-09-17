package callback

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type Context struct {
	Payload    io.Reader
	RouterPath string
}

func NewContext[T any](data T, routerPath string) *Context {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling context: %v", err)
		return nil
	}

	return &Context{Payload: bytes.NewReader(payload), RouterPath: routerPath}
}
