package controller

import (
	"example/web-service-gin/controller/user"
	"example/web-service-gin/repository"
)

type Controllers struct {
	userController *user.Controller
}

func InitControllers(repositories *repository.Repositories) *Controllers {
	return &Controllers{
		userController: user.InitController(repositories.UserRepo),
	}
}
