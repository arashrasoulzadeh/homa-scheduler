package invokers

import (
	"context"
	"net"
	"os"

	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"github.com/arashrasoulzadeh/homa-scheduler/socket"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "3001"
	SERVER_TYPE = "tcp"
)

func RunSocketServer(lc fx.Lifecycle, logger *zap.SugaredLogger, data providers.Data) *net.Listener {
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err != nil {
				logger.Errorw("Error listening", "details", err.Error())
				os.Exit(1)
			}
			logger.Infow("socket running", "Host", SERVER_HOST, "PORT", SERVER_PORT)
			go socket.Init(server, logger, data)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Close()
			return nil
		},
	})

	return &server
}
