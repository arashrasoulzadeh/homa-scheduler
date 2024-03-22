package models

import (
	"time"

	"github.com/google/uuid"
)

type CommandDelivery struct {
	ID          uuid.UUID
	CommandId   uuid.UUID `gorm:"size:40"`
	Command     Command
	ClientId    uuid.UUID `gorm:"size:40"`
	Client      Clinet
	InstanceId  uuid.UUID `gorm:"size:40"`
	Instance    Instance
	Room        string
	DeliveredAt time.Time
}
