package repository

import (
	"example/web-service-gin/repository/rolerepo"
	"example/web-service-gin/repository/userrepo"

	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo *userrepo.UserRepo
	RoleRepo *rolerepo.RoleRepo
}

func InitRepositories(db *gorm.DB) *Repositories {
	userrepo := userrepo.NewUserRepo(db)
	rolerepo := rolerepo.NewRoleRepo(db)
	return &Repositories{
		UserRepo: userrepo,
		RoleRepo: rolerepo,
	}
}
