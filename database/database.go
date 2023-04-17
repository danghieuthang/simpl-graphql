package database

import (
	"os"

	"example/web-service-gin/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitializeDatabase() {
	var err error
	psqlConn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(psqlConn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
}

func SyncDatabase() {
	DB.Migrator().AutoMigrate(&entity.User{}, &entity.Role{})
}
