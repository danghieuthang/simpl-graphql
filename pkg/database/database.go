package database

import (
	"example/web-service-gin/pkg/entity"
	"os"

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

	DB.AutoMigrate(&entity.Role{}, &entity.User{})
}
