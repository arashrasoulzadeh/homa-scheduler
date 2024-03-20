package invokers

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"slices"
	"strings"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunSocketServer(lc fx.Lifecycle, logger *zap.SugaredLogger, data providers.Data) *net.Listener {

	server := providers.SocketServer()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go listen(server, logger, data)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Close()
			return nil
		},
	})

	return &server
}

func listen(server net.Listener, logger *zap.SugaredLogger, data providers.Data) {
	for {
		connection, err := server.Accept()
		if err != nil {
			panic(err)
		}
		PorcessClient(connection, logger, data)
	}
}

func PorcessClient(connection net.Conn, logger *zap.SugaredLogger, data providers.Data) {
	clinet := &models.Clinet{
		Connection: connection,
		Channels:   []string{},
	}
	logger.Infow("client connected", "addr", connection.RemoteAddr())
	go readBus(data, clinet, logger)
	go ProcessMessage(clinet, data.Connection, logger)
}
func readBus(data providers.Data, client *models.Clinet, logger *zap.SugaredLogger) {
	for {
		bus := <-data.Bus
		command := bus.(models.Command)
		if slices.Contains(client.Channels, command.Channel) {
			sendMesasge(client, createMessage("new command", command.Id, 200), logger)
		}
	}
}

func ProcessMessage(client *models.Clinet, db *gorm.DB, logger *zap.SugaredLogger) {
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

func sendMesasge(client *models.Clinet, payload []byte, logger *zap.SugaredLogger) {
	client.Connection.Write(payload)
	logger.Infow("outgoing", "addr", client.Connection.RemoteAddr(), "payload", string(payload))
}

func execCommand(cmd string, client *models.Clinet, db *gorm.DB, logger *zap.SugaredLogger) error {
	parts := strings.Split(cmd, ":")
	if len(parts) == 2 {
		cmd := parts[0]
		args := strings.ReplaceAll(parts[1], "\r\n", "")
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
