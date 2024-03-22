package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Command struct {
	ID        uuid.UUID
	Command   string
	Args      datatypes.JSONMap
	Tags      string
	Channel   string
	Status    string
	CreatedAt time.Time
}

func (c *Command) MarkAsDev() {
	c.AddTag("dev")
}

func (c *Command) MarkAsInPending() {
	c.Status = "Pending"
}

func (c *Command) MarkAsInProgress() {
	c.Status = "In Progress"
}

func (c *Command) AddTag(tag string) {
	if c.Tags == "" {
		c.Tags = tag
	} else {
		c.Tags = c.Tags + "," + tag
	}
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
func (c *Command) TagsArray() []string {
	return strings.Split(c.Tags, ",")
}
