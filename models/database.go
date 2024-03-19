package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := "root:root@tcp(localhost:8889)/homa-dispatcher?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	fmt.Println("connected!")
	return db
}

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(Command{})
}
