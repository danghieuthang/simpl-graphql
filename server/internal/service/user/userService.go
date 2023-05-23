package user

import (
	"context"
	"errors"
	"example/web-service-gin/internal/constant/app_error"
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	View(id int, fields []string) (*entity.User, error)
	List(name string, page int, pageSize int) (*[]entity.User, int, error)
	Update(u *entity.User, ctx context.Context) (*entity.User, error)
	Create(u *entity.User, ctx context.Context) (*entity.User, error)
	Login(email string, password string) (*entity.User, error)
}

type userService struct {
	repository repository.IRepository
}

func InitUserService(repository repository.IRepository) IUserService {
	return &userService{
		repository: repository,
	}
}

func (c *userService) View(id int, fields []string) (*entity.User, error) {
	var user entity.User
	condition := fmt.Sprintf("id = %d ", id)
	err := c.repository.GetWhere(&user, condition, "Role")
	return &user, err
}

func (c *userService) Create(u *entity.User, ctx context.Context) (*entity.User, error) {
	var user entity.User
	condition := fmt.Sprintf("id = %d or email = '%s'", u.Id, u.Email)
	err := c.repository.GetWhere(&user, condition)
	if user.Id == u.Id {
		return nil, errors.New(app_error.USER_ID_EXIST)
	}
	if user.Email == u.Email {
		return nil, errors.New(app_error.USER_EMAIL_EXIST)
	}
	err = c.repository.CreateWithContext(u, ctx)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c *userService) Update(u *entity.User, ctx context.Context) (*entity.User, error) {
	err := c.repository.SaveWithContext(u, ctx)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *userService) List(name string, page int, pageSize int) (*[]entity.User, int, error) {
	var users []entity.User
	var total int64
	condition := fmt.Sprintf("name ilike '%%%v%%' ", name)
	c.repository.GetWhereBatch(&users, condition, pageSize, (page-1)*pageSize)
	c.repository.CountWhere(&users, &total, condition)
	return &users, int(total), nil
}

func (c *userService) Login(email string, password string) (*entity.User, error) {
	var user entity.User
	condition := fmt.Sprintf("email = '%s' ", email)
	err := c.repository.GetWhere(&user, condition)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, errors.New("User not exist")
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Password incorrect")
	}
	return &user, err
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
