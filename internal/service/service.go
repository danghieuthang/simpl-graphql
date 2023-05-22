package service

import (
	"example/web-service-gin/internal/service/role"
	"example/web-service-gin/internal/service/user"
	"example/web-service-gin/pkg/repository"
)

type ServiceFactory struct {
	UserService user.IUserService
	RoleService role.IRoleService
}

func InitServices(repository repository.IRepository) *ServiceFactory {
	return &ServiceFactory{
		UserService: user.InitUserService(repository),
		RoleService: role.InitRoleService(repository),
	}
}
