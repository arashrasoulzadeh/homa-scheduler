package models

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(logger *zap.SugaredLogger) *gorm.DB {
	dsn := "root:root@tcp(localhost:8889)/homa-dispatcher?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	logger.Info("connected to database")
	return db
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(Command{})
}
