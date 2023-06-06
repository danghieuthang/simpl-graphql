package resolver

import (
	"context"
	"example/web-service-gin/graph/model"
	"example/web-service-gin/pkg/utils"
)

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, page *model.PageArgs) (*model.RolePage, error) {
	roles, total, _ := r.ServiceFactory.RoleService.List(ctx, "", page.Page, page.PageSize)
	var roleViews []*model.Role

	for _, v := range roles {
		role := utils.Map[model.Role](v)
		roleViews = append(roleViews, role)
	}
	return &model.RolePage{
		Data: roleViews,
		PageInfo: &model.PageInfo{
			Total:    total,
			Page:     page.Page,
			PageSize: page.PageSize,
		},
	}, nil
}
