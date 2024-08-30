package main

import (
	"encoding/gob"
	"net"
)

type socket struct {
	addr string
}

func (socket socket) dial() (net.Conn, error) {
	return net.Dial("unix", socket.addr)
}

func (socket socket) send(encoder *gob.Encoder, v string) error {
	return encoder.Encode([]byte(v + "\n"))
}

func (socket socket) receive(decoder *gob.Decoder) (string, error) {
	var data []byte
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}
	return string(data), nil
}
