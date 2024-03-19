package main

import (
	"fmt"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/arashrasoulzadeh/homa-scheduler/router"
	"github.com/arashrasoulzadeh/homa-scheduler/socket"
)

func main() {
	fmt.Println("hello world")
	connection := models.Connect()
	models.RunMigrations(connection)
	go router.Init(connection)
	go socket.Init(connection)
	for true {

	}
}
