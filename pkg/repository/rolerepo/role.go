package rolerepo

import (
	"errors"
	"example/web-service-gin/pkg/entity"

	"gorm.io/gorm"
)

type IRoleRepo interface {
	Create(r *entity.Role) (*entity.Role, error)
}

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) IRoleRepo {
	return &RoleRepo{
		db: db,
	}
}
func (repo *RoleRepo) Create(r *entity.Role) (*entity.Role, error) {
	var isExistId int
	repo.db.Raw("SELECT id FROM roles where id=?", *&r.Id).Scan(&isExistId)
	if isExistId == 0 {
		repo.db.Create(&r)
		return r, nil
	}
	return nil, errors.New("Id was exist")
}
