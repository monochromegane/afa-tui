package main

import (
	"encoding/gob"
	"net"
	"time"
)

type socket struct{}

func (socket socket) dial() (net.Conn, error) {
	time.Sleep(3 * time.Second)
	return nil, nil
}

func (socket socket) send(encoder *gob.Encoder, v string) error {
	time.Sleep(3 * time.Second)
	return nil
}

func (socket socket) receive(decoder *gob.Decoder) (string, error) {
	time.Sleep(3 * time.Second)
	return "", nil
}
