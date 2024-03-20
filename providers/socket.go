package providers

import "net"

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "3001"
	SERVER_TYPE = "tcp"
)

func SocketServer() net.Listener {
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		panic(err)
	}

	return server
}
