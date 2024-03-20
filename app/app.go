package app

import (
	"github.com/arashrasoulzadeh/homa-scheduler/invokers"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"go.uber.org/fx"
)

func Run() {
	fx.New(
		fx.Provide(providers.LoggingSugar, providers.ConnectToDatabase),
		fx.Invoke(invokers.RunHttpServer, invokers.RunSocketServer),
	).Run()
}
