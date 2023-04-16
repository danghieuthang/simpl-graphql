package user

import (
	"example/web-service-gin/entity"
	"example/web-service-gin/repository/userrepo"
)

type Controller struct {
	service userrepo.IUserRepo
}

func InitController(userRepo *userrepo.UserRepo) *Controller {
	return &Controller{
		service: userRepo,
	}
}

func (c *Controller) View(id int) (*entity.User, error) {
	user, err := c.service.View(id)
	return user, err
}

func (c *Controller) Create(u *entity.User) (*entity.User, error) {
	res, err := c.service.Create(u)
	return res, err
}

func (c *Controller) Update(u *entity.User) (*entity.User, error) {
	res, err := c.service.Update(u)
	return res, err
}

func (c *Controller) List(name string) (*[]entity.User, error) {
	res, err := c.service.List(name)
	return res, err
}

func (c *Controller) Login(email string, password string) (*entity.User, error) {
	res, err := c.service.Login(email, password)
	return res, err
}
