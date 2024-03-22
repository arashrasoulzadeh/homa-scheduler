package invokers

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunSocketServer(lc fx.Lifecycle, logger *zap.SugaredLogger, data providers.Data) *net.Listener {

	server := providers.SocketServer()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			hostname, _ := os.Hostname()

			instance := models.Instance{
				ID:        uuid.New(),
				LocalAddr: GetLocalIP(),
				HostName:  hostname,
			}
			data.Connection.Save(instance)
			go listen(server, logger, data, instance)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Close()
			return nil
		},
	})

	return &server
}

func listen(server net.Listener, logger *zap.SugaredLogger, data providers.Data, instance models.Instance) {
	for {
		connection, err := server.Accept()
		if err != nil {
			panic(err)
		}
		PorcessClient(connection, logger, data, instance)
	}
}

func PorcessClient(connection net.Conn, logger *zap.SugaredLogger, data providers.Data, instance models.Instance) {
	client := &models.Clinet{
		ID:       uuid.New(),
		Channels: []string{},
		Instance: instance,
	}
	client.Persist(data.Connection)
	client.SetConnection(&connection)
	logger.Infow("client connected", "addr", connection.RemoteAddr())
	go readBus(data, client, logger)
	go ProcessMessage(client, data.Connection, logger)
}
func readBus(data providers.Data, client *models.Clinet, logger *zap.SugaredLogger) {
	for {
		bus := <-data.Bus
		command := bus.(models.Command)
		if slices.Contains(client.Channels, command.Channel) {
			sendMesasge(client, createMessage("new command", command.ID, 200, command.Channel), logger)
		}
	}
}

func ProcessMessage(client *models.Clinet, db *gorm.DB, logger *zap.SugaredLogger) {
	defer client.GetConnection().Close()
	for {
		buffer := make([]byte, 1024)
		mLen, err := client.GetConnection().Read(buffer)
		if err != nil {
			logger.Error("error reading incoming payload", "addr", client.GetConnection().RemoteAddr(), "error", err.Error())
		}
		logger.Infow("incoming", "addr", client.GetConnection().RemoteAddr(), "payload", string(buffer[:mLen]))

		err = execCommand(string(buffer[:mLen]), client, db, logger)
		if err != nil {
			logger.Error("error executing command", "addr", client.GetConnection().RemoteAddr(), "error", err.Error())
			sendMesasge(client, []byte("error:"+err.Error()+"\n"), logger)
		}
	}

}

func sendMesasge(client *models.Clinet, payload []byte, logger *zap.SugaredLogger) {
	client.GetConnection().Write(payload)
	logger.Infow("outgoing", "addr", client.GetConnection().RemoteAddr(), "payload", string(payload))
}

func execCommand(cmd string, client *models.Clinet, db *gorm.DB, logger *zap.SugaredLogger) error {
	parts := strings.Split(cmd, ":")
	if len(parts) == 2 {
		cmd := parts[0]
		args := strings.ReplaceAll(parts[1], "\r\n", "")
		switch cmd {
		case "join":
			client.AddChannel(args)
			sendMesasge(client, createMessage("joined channel!", client.Channels, 200, "general"), logger)
		case "get":
			getNewCommand(db, client, logger)
		}
		return nil
	}
	sendMesasge(client, createMessage("cannot process command!", nil, 500, "general"), logger)
	return errors.New("cannot process command")
}

func getNewCommand(db *gorm.DB, client *models.Clinet, logger *zap.SugaredLogger) {
	var cmd *models.Command
	db.First(&cmd)
	if cmd.Status == "" {
		sendMesasge(client, createMessage("can't get command!", nil, 500, "general"), logger)
		return
	}
	cmd.MarkAsInProgress()
	db.Save(cmd)
	delivery := models.CommandDelivery{
		ID:          uuid.New(),
		Command:     *cmd,
		Client:      *client,
		Instance:    client.Instance,
		DeliveredAt: time.Now(),
	}
	db.Save(&delivery)
	sendMesasge(client, createMessage("new command!", cmd, 200, "general"), logger)
}

func createMessage(message string, data interface{}, status int, channel string) []byte {
	msg := models.Message{
		Message: message,
		Data:    data,
		Status:  status,
		Channel: channel,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return []byte("{}")
	}
	return b
}
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
