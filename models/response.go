package models

type Message struct {
	Message string
	Data    interface{}
	Status  int
	Channel string
}
