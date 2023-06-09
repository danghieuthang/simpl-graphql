package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"example/web-service-gin/graph/model"
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/middleware/auth"
	"example/web-service-gin/pkg/utils"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.User, error) {
	user, err := r.ServiceFactory.UserService.Create(ctx, &entity.User{
		Id:       input.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
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
	user, err := r.ServiceFactory.UserService.Update(ctx, &entity.User{
		Id:       input.ID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
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
	res, err := r.ServiceFactory.UserService.Login(ctx, input.Email, input.Password)
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
	preloads := utils.GetPreloads(ctx)
	user, err := r.ServiceFactory.UserService.View(ctx, id, preloads)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Id != id {
		return nil, errors.New("Id does not exist")
	}

	viewUser := &model.User{}
	ms, _ := json.Marshal(user)
	json.Unmarshal(ms, &viewUser)
	return viewUser, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, page *model.PageArgs) (*model.UserPage, error) {
	users, total, _ := r.ServiceFactory.UserService.List(ctx, "", page.Page, page.PageSize)

	var viewUsers []*model.User

	for _, v := range users {
		user := utils.Map[model.User](v)
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
	preloads := utils.GetPreloads(ctx)
	domainUser, _ := r.ServiceFactory.UserService.View(ctx, authenUser.Id, preloads)
	viewUser := &model.User{}
	ms, _ := json.Marshal(domainUser)
	json.Unmarshal(ms, &viewUser)
	return viewUser, nil
}
