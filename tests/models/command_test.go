package tests

import (
	"slices"
	"testing"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func sampleCommand() models.Command {
	return models.Command{
		Id:      uuid.New(),
		Command: "hello",
		Args:    datatypes.JSONMap{},
		Tags:    "",
	}
}
func TestEmptyId(t *testing.T) {
	command := sampleCommand()
	if command.Id == uuid.Nil {
		t.Fail()
	}
}
func TestMarkAsDev(t *testing.T) {
	command := sampleCommand()
	command.MarkAsDev()
	if !slices.Contains(command.TagsArray(), "dev") {
		t.Fail()
	}
}
func TestSerialize(t *testing.T) {
	command := sampleCommand()
	_, err := command.Serialize()
	if err != nil {
		t.Fail()
	}
}
func TestDeSerialize(t *testing.T) {
	command := sampleCommand()
	data, err := command.Serialize()
	if err != nil {
		t.Fail()
	}
	command.Deserialize(data)
	if command.Command != "hello" {
		t.Fail()
	}
}
