package providers

import (
	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Data struct {
	Connection *gorm.DB
	Bus        chan interface{}
}

func ConnectToDatabase(logger *zap.SugaredLogger) Data {
	return Data{Connection: models.Connect(logger), Bus: make(chan interface{})}
}
