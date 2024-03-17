package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Command struct {
	Id      uuid.UUID
	Command string
	Args    map[string]string
	Tags    []string
}

func (c *Command) MarkAsDev() {
	c.Tags = append(c.Tags, "dev")
}

func (c *Command) Serialize() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}
func (c *Command) Deserialize(data []byte) error {
	err := json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}
