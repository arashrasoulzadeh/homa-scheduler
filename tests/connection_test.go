package test

import (
	"net"
	"testing"
	"time"

	"github.com/arashrasoulzadeh/homa-scheduler/providers"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

func TestSocketServer(t *testing.T) {
	server := providers.SocketServer()
	server.Accept()
	_, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	time.Sleep(5 * time.Second)
	if err != nil {
		t.Fail()
	}
}
