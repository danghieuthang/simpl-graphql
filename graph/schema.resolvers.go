package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.31

import (
	"context"
	"encoding/json"
	"example/web-service-gin/graph/model"
	"example/web-service-gin/internal/auth"
	"example/web-service-gin/pkg/entity"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user, err := r.ServiceFactory.UserService.Create(&entity.User{
		Id:       input.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}, ctx)
	if err != nil {
		return nil, err
	}
	viewUser := &model.User{}
	ms, _ := json.Marshal(user)
	json.Unmarshal(ms, &user)
	return viewUser, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user, err := r.ServiceFactory.UserService.Update(&entity.User{
		Id:       input.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}, ctx)
	if err != nil {
		return nil, err
	}
	viewUser := &model.User{}
	ms, _ := json.Marshal(user)
	json.Unmarshal(ms, &user)
	return viewUser, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthenType, error) {
	res, err := r.ServiceFactory.UserService.Login(input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	jwtToken, err := auth.GenerateToken(res)
	return &model.AuthenType{
		Token:     jwtToken,
		TokenType: "jwt",
		ExpiresIn: 60 * 60,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	// panic(fmt.Errorf("not implemented: User - user"))
	preloads := GetPreloads(ctx)
	user, err := r.ServiceFactory.UserService.View(id, preloads)
	if err != nil {
		return nil, err
	}

	viewUser := &model.User{}
	ms, _ := json.Marshal(user)
	json.Unmarshal(ms, &viewUser)
	return viewUser, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, page *model.PageArgs) (*model.UserPage, error) {
	users, total, _ := r.ServiceFactory.UserService.List("", page.Page, page.PageSize)
	var viewUsers []*model.User
	for _, v := range *users {
		user := &model.User{}
		ms, _ := json.Marshal(v)
		json.Unmarshal(ms, &user)
		viewUsers = append(viewUsers, user)
	}
	return &model.UserPage{
		Data: viewUsers,
		PageInfo: &model.PageInfo{
			Total:    total,
			Page:     page.Page,
			PageSize: page.PageSize,
		},
	}, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	authenUser := ctx.Value("currentUser").(*auth.AuthenticatedUser)
	preloads := GetPreloads(ctx)
	domainUser, _ := r.ServiceFactory.UserService.View(authenUser.Id, preloads)
	viewUser := &model.User{}
	ms, _ := json.Marshal(domainUser)
	json.Unmarshal(ms, &viewUser)
	return viewUser, nil
}

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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
