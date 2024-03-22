package models

import "github.com/google/uuid"

type Instance struct {
	ID        uuid.UUID
	LocalAddr string
	HostName  string
}
