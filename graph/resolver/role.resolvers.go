package resolver

import (
	"context"
	"encoding/json"
	"example/web-service-gin/graph/model"
)

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, page *model.PageArgs) (*model.RolePage, error) {
	roles, total, _ := r.ServiceFactory.RoleService.List("", page.Page, page.PageSize)
	var roleViews []*model.Role
	for _, v := range *roles {
		role := &model.Role{}
		ms, _ := json.Marshal(v)
		json.Unmarshal(ms, &role)
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
