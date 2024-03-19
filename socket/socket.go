package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"gorm.io/gorm"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "3001"
	SERVER_TYPE = "tcp"
)

func Init(db *gorm.DB) {
	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection, db)
	}
}

func processClient(connection net.Conn, db *gorm.DB) {
	defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Println("Received: ", string(buffer[:mLen]))
		err = execCommand(string(buffer[:mLen]), connection, db)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			connection.Write([]byte("error:" + err.Error() + "\n"))

		}
	}

}

func execCommand(cmd string, connection net.Conn, db *gorm.DB) error {
	parts := strings.Split(cmd, ":")
	if len(parts) == 2 {
		cmd := parts[0]
		// args := parts[1]
		switch cmd {
		case "get":

			var cmd *models.Command
			db.First(&cmd)
			if cmd == nil {
				connection.Write(createMessage("can't get command!", nil, 500))
				break
			}
			cmd.MarkAsInProgress()
			db.Save(cmd)

			connection.Write(createMessage("new command!", cmd, 200))
		}
		return nil
	}
	connection.Write(createMessage("cannot process command!", nil, 500))
	return errors.New("cannot process command!")
}

func createMessage(message string, data interface{}, status int) []byte {
	msg := models.Message{
		Message: message,
		Data:    data,
		Status:  status,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return []byte("{}")
	}
	return b
}
