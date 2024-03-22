package models

import (
	"net"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Clinet struct {
	ID         uuid.UUID
	connection net.Conn
	Channels   datatypes.JSONSlice[string]
	InstanceId uuid.UUID `gorm:"size:40"`
	Instance   Instance
}

func (client *Clinet) AddChannel(channel string) {
	client.Channels = append(client.Channels, channel)
}

func (client *Clinet) SetConnection(connection *net.Conn) *net.Conn {
	client.connection = *connection
	return &client.connection
}
func (client *Clinet) GetConnection() net.Conn {
	return client.connection
}

func (client *Clinet) Persist(db *gorm.DB) {
	db.Create(client)
}
