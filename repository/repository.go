package repository

import (
	"example/web-service-gin/repository/rolerepo"
	"example/web-service-gin/repository/userrepo"

	"gorm.io/gorm"
)

type RepositoryFactory struct {
	UserRepo userrepo.IUserRepo
	RoleRepo rolerepo.IRoleRepo
}

func InitRepositories(db *gorm.DB) *RepositoryFactory {
	userrepo := userrepo.NewUserRepo(db)
	rolerepo := rolerepo.NewRoleRepo(db)
	return &RepositoryFactory{
		UserRepo: userrepo,
		RoleRepo: rolerepo,
	}
}
