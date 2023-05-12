package user

import (
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
	Update(u *entity.User) (*entity.User, error)
	Create(u *entity.User) (*entity.User, error)
	Login(email string, password string) (*entity.User, error)
}

type UserService struct {
	repository repository.IRepository
}

func InitUserService(repository repository.IRepository) IUserService {
	return &UserService{
		repository: repository,
	}
}

func (c *UserService) View(id int, fields []string) (*entity.User, error) {
	var user entity.User
	condition := fmt.Sprintf("id = %d ", id)
	err := c.repository.GetWhere(&user, condition, "Role")
	return &user, err
}

func (c *UserService) Create(u *entity.User) (*entity.User, error) {
	var user entity.User
	condition := fmt.Sprintf("id = %d or email = '%s'", u.Id, u.Email)
	err := c.repository.GetWhere(&user, condition)
	if user.Id == u.Id {
		return nil, errors.New(app_error.USER_ID_EXIST)
	}
	if user.Email == u.Email {
		return nil, errors.New(app_error.USER_EMAIL_EXIST)
	}
	err = c.repository.Create(u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c *UserService) Update(u *entity.User) (*entity.User, error) {
	err := c.repository.Save(u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c *UserService) List(name string, page int, pageSize int) (*[]entity.User, int, error) {
	var users []entity.User
	var total int64
	condition := fmt.Sprintf("name ilike '%%%v%%' ", name)
	c.repository.GetWhereBatch(&users, condition, pageSize, (page-1)*pageSize)
	c.repository.CountWhere(&users, &total, condition)
	return &users, int(total), nil
}

func (c *UserService) Login(email string, password string) (*entity.User, error) {
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
