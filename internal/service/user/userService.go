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
	View(ctx context.Context, id int, fields []string) (*entity.User, error)
	List(ctx context.Context, name string, page int, pageSize int) ([]*entity.User, int, error)
	Update(ctx context.Context, u *entity.User) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) (*entity.User, error)
	Login(ctx context.Context, email string, password string) (*entity.User, error)
}

type userService struct {
	repository repository.IBaseRepository[*entity.User]
}

func InitUserService(repository repository.IBaseRepository[*entity.User]) IUserService {
	return &userService{
		repository: repository,
	}
}

func (c *userService) View(ctx context.Context, id int, fields []string) (*entity.User, error) {
	var user *entity.User
	condition := fmt.Sprintf("id = %d ", id)
	user, err := c.repository.GetOneAsNoTracking(ctx, condition, "Role")
	return user, err
}

func (c *userService) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	var user *entity.User
	condition := fmt.Sprintf("id = %d or email = '%s'", u.Id, u.Email)
	user, err := c.repository.GetOneAsNoTracking(ctx, condition)
	if user.Id == u.Id {
		return nil, errors.New(app_error.USER_ID_EXIST)
	}
	if user.Email == u.Email {
		return nil, errors.New(app_error.USER_EMAIL_EXIST)
	}
	u, err = c.repository.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (c *userService) Update(ctx context.Context, u *entity.User) (*entity.User, error) {
	domainUser, err := c.repository.GetOne(ctx, fmt.Sprintf("id=%v", u.Id))
	if err != nil || domainUser == nil || domainUser.Id != u.Id {
		return nil, errors.New("User invalid")
	}
	err = c.repository.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *userService) List(ctx context.Context, name string, page int, pageSize int) ([]*entity.User, int, error) {
	var total int64
	condition := fmt.Sprintf("name ilike '%%%v%%' ", name)
	total, err := c.repository.Count(ctx, condition)
	if total == 0 {
		return nil, 0, nil
	}
	users, err := c.repository.GetAllAsNoTracking(ctx, condition, pageSize, (page-1)*pageSize)
	return users, int(total), err
}

func (c *userService) Login(ctx context.Context, email string, password string) (*entity.User, error) {
	condition := fmt.Sprintf("email = '%s' ", email)
	user, err := c.repository.GetOneAsNoTracking(ctx, condition, "Role")
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, errors.New("User not exist")
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Password incorrect")
	}
	return user, err
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
