package socket

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Clinet struct {
	Connection net.Conn
	Channels   []string
}

func Init(server net.Listener, logger *zap.SugaredLogger, data providers.Data) {
	for {
		connection, err := server.Accept()
		clinet := Clinet{
			Connection: connection,
			Channels:   []string{},
		}
		if err != nil {
			logger.Error("socket running", "Error", err.Error())
			os.Exit(1)
		}
		logger.Infow("client connected", "addr", connection.RemoteAddr())
		go readBus(data, clinet, logger)
		go processClient(clinet, data.Connection, logger)
	}
}

func readBus(data providers.Data, client Clinet, logger *zap.SugaredLogger) {
	for {
		bus := <-data.Bus
		command := bus.(models.Command)
		if slices.Contains(client.Channels, command.Channel) {
			sendMesasge(client, createMessage("new command", command.Id, 200), logger)
		}
	}
}

func processClient(client Clinet, db *gorm.DB, logger *zap.SugaredLogger) {
	defer client.Connection.Close()
	for {
		buffer := make([]byte, 1024)
		mLen, err := client.Connection.Read(buffer)
		if err != nil {
			logger.Error("error reading incoming payload", "addr", client.Connection.RemoteAddr(), "error", err.Error())
		}
		logger.Infow("incoming", "addr", client.Connection.RemoteAddr(), "payload", string(buffer[:mLen]))

		err = execCommand(string(buffer[:mLen]), client, db, logger)
		if err != nil {
			logger.Error("error executing command", "addr", client.Connection.RemoteAddr(), "error", err.Error())
			sendMesasge(client, []byte("error:"+err.Error()+"\n"), logger)
		}
	}

}

func sendMesasge(client Clinet, payload []byte, logger *zap.SugaredLogger) {
	client.Connection.Write(payload)
	logger.Infow("outgoing", "addr", client.Connection.RemoteAddr(), "payload", string(payload))
}

func execCommand(cmd string, client Clinet, db *gorm.DB, logger *zap.SugaredLogger) error {
	parts := strings.Split(cmd, ":")
	if len(parts) == 2 {
		cmd := parts[0]
		args := parts[1]
		switch cmd {
		case "join":
			client.Channels = append(client.Channels, args)
			sendMesasge(client, createMessage("joined channel!", client.Channels, 500), logger)
			break
		case "get":
			var cmd *models.Command
			db.First(&cmd)
			if cmd == nil {
				sendMesasge(client, createMessage("can't get command!", nil, 500), logger)
				break
			}
			cmd.MarkAsInProgress()
			db.Save(cmd)
			sendMesasge(client, createMessage("new command!", cmd, 200), logger)
		}
		return nil
	}
	sendMesasge(client, createMessage("cannot process command!", nil, 500), logger)
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
