package models

import "net"

type Clinet struct {
	Connection net.Conn
	Channels   []string
}
