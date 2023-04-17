package controller

import (
	"example/web-service-gin/controller/user"
	"example/web-service-gin/repository"
)

type ControllerFactory struct {
	UserController *user.Controller
}

func InitControllers(repositories *repository.RepositoryFactory) *ControllerFactory {
	return &ControllerFactory{
		UserController: user.InitController(repositories.UserRepo),
	}
}
