package userrepo

import (
	"errors"
	"example/web-service-gin/internal/constant/app_error"
	"example/web-service-gin/pkg/entity"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepo interface {
	View(id int, fields []string) (*entity.User, error)
	List(name string) (*[]entity.User, error)
	Update(u *entity.User) (*entity.User, error)
	Create(u *entity.User) (*entity.User, error)
	Login(email string, password string) (*entity.User, error)
	// Delete(*entity.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	return &UserRepo{
		db: db,
	}
}
func (repo *UserRepo) View(id int, fields []string) (*entity.User, error) {
	var user entity.User
	repo.db.Where("users.id=?", id).Preload(clause.Associations).Find(&user)
	// utils.PreLoadWithSelect(repo.db.Debug(), fields, "user").Where("users.id=?", id).First(&user)
	return &user, nil
}

func (repo *UserRepo) Login(email string, password string) (*entity.User, error) {
	var user entity.User
	repo.db.Where("email=?", email).Preload(clause.Associations).Find(&user)
	if user.Id == 0 {
		return nil, errors.New("User not exist")
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Password incorrect")
	}
	return &user, nil
}

func (repo *UserRepo) List(name string) (*[]entity.User, error) {
	var users []entity.User
	if len(name) == 0 {
		repo.db.Find(&users)
	} else {
		repo.db.Where("name ilike ?", "%"+name+"%").Find(&users)
	}

	return &users, nil
}

func (repo *UserRepo) Create(u *entity.User) (*entity.User, error) {
	u.Password, _ = HashPassword(u.Password)
	var isExistId int
	repo.db.Raw("SELECT id FROM users where id=?", *&u.Id).Scan(&isExistId)
	if isExistId == 0 {
		repo.db.Create(&u)
		return u, nil
	}
	return nil, errors.New(app_error.USER_ID_EXIST_ID)
}

func (repo *UserRepo) Update(u *entity.User) (*entity.User, error) {
	var dbUser entity.User
	repo.db.Where("id=?", &u.Id).First(&dbUser)
	if dbUser.Id == 0 {

		return nil, errors.New("Id not exist")
	}

	dbUser.Name = u.Name
	dbUser.Email = u.Email
	dbUser.LastModifiedAt = time.Now()
	repo.db.Save(&dbUser)

	return &dbUser, nil
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
