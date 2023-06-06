package service

import (
	"example/web-service-gin/internal/service/role"
	"example/web-service-gin/internal/service/user"
	"example/web-service-gin/pkg/logger"
	"example/web-service-gin/pkg/repository"

	"gorm.io/gorm"
)

type ServiceFactory struct {
	UserService user.IUserService
	RoleService role.IRoleService
}

func InitServices(db *gorm.DB, logger logger.ILogger) *ServiceFactory {
	return &ServiceFactory{
		UserService: user.InitUserService(repository.NewUserRepository(db, logger)),
		RoleService: role.InitRoleService(repository.NewRoleRepository(db, logger)),
	}
}
