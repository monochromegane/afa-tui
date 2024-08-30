package main

import "net"

type errMsg struct{ error }
type connectedMsg struct{ conn net.Conn }
type promptMsg struct{}
type sentMsg struct{}
type closeMsg struct{ error }
type responseMsg struct{ message string }
