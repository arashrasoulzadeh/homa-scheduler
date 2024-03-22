package models

import "github.com/google/uuid"

type Room struct {
	Id   uuid.UUID
	name string
}
