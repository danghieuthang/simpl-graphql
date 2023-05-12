package role

import (
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/repository"
	"fmt"
)

type IRoleService interface {
	View(id int, fields []string) (*entity.Role, error)
	List(name string, page int, pageSize int) (*[]entity.Role, int, error)
}

type RoleService struct {
	repository repository.IRepository
}

func (c *RoleService) List(name string, page int, pageSize int) (*[]entity.Role, int, error) {
	var roles []entity.Role
	var total int64
	condition := fmt.Sprintf("name ilike '%%%v%%' ", name)
	c.repository.GetWhereBatch(&roles, condition, pageSize, (page-1)*pageSize)
	c.repository.CountWhere(&roles, &total, condition)
	return &roles, int(total), nil
}
func (c *RoleService) View(id int, fields []string) (*entity.Role, error) {
	var user entity.Role
	c.repository.GetOneByField(&user, "id", id)
	return &user, nil
}
func InitRoleService(repository repository.IRepository) IRoleService {
	return &RoleService{
		repository: repository,
	}
}
