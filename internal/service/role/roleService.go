package role

import (
	"context"
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/repository"
	"fmt"
)

type IRoleService interface {
	View(ctx context.Context, id int, fields []string) (*entity.Role, error)
	List(ctx context.Context, name string, page int, pageSize int) ([]*entity.Role, int, error)
}

type roleService struct {
	repository repository.IBaseRepository[*entity.Role]
}

func (c *roleService) List(ctx context.Context, name string, page int, pageSize int) ([]*entity.Role, int, error) {
	condition := fmt.Sprintf("name ilike '%%%v%%' ", name)
	total, err := c.repository.Count(ctx, condition)
	if total == 0 {
		return nil, 0, nil
	}
	roles, err := c.repository.GetAllAsNoTracking(ctx, condition, pageSize, (page-1)*pageSize)
	return roles, int(total), err
}
func (c *roleService) View(ctx context.Context, id int, fields []string) (*entity.Role, error) {
	role, err := c.repository.GetOneAsNoTracking(ctx, fmt.Sprintf("id=?", id))
	return role, err
}
func InitRoleService(repository repository.IBaseRepository[*entity.Role]) IRoleService {
	return &roleService{
		repository: repository,
	}
}
