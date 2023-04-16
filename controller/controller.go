package controller

import (
	"example/web-service-gin/controller/user"
	"example/web-service-gin/repository"
)

type ControllerFactory struct {
	userController *user.Controller
}

func InitControllers(repositories *repository.RepositoryFactory) *ControllerFactory {
	return &ControllerFactory{
		userController: user.InitController(repositories.UserRepo),
	}
}
