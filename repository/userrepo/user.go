package userrepo

import (
	"errors"
	"example/web-service-gin/entity"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IUserRepo interface {
	View(id int) (*entity.User, error)
	List(name string) (*[]entity.User, error)
	Update(u *entity.User) (*entity.User, error)
	Create(u *entity.User) (*entity.User, error)
	// Delete(*entity.User) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo *UserRepo) View(id int) (*entity.User, error) {
	var user entity.User
	repo.db.Preload(clause.Associations).Where("users.id=?", id).First(&user)
	if user.Id == 0 {
		return nil, errors.New("Id not exist")
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
	var isExistId int
	repo.db.Raw("SELECT id FROM users where id=?", *&u.Id).Scan(&isExistId)
	if isExistId == 0 {
		repo.db.Create(&u)
		return u, nil
	}
	return nil, errors.New("Id was exist")

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
