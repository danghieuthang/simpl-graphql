package database

import (
	"fmt"

	"example/web-service-gin/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "golang"
)

func GetDatabase() *gorm.DB {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.Migrator().AutoMigrate(&entity.User{}, &entity.Role{})
	return db
}
