package main

import (
	"github.com/arashrasoulzadeh/homa-scheduler/invokers"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(providers.LoggingSugar, providers.ConnectToDatabase),
		fx.Invoke(invokers.RunHttpServer, invokers.RunSocketServer),
	).Run()

	// connection := models.Connect(sugar)
	// models.RunMigrations(connection)
	// for true {

	// }
}
