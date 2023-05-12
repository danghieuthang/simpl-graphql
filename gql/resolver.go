package gql

import (
	"example/web-service-gin/controller"
	"example/web-service-gin/entity"
	"example/web-service-gin/pkg/jwt"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

type Resolver struct {
	ServiceFactory *controller.ServiceFactory
}

func (r *Resolver) CreateUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	password, _ := params.Args["password"].(string)
	res, err := r.ServiceFactory.UserService.Create(&entity.User{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, gqlerrors.FormatError(err)
	}
	return res, nil
}

func (r *Resolver) Login(params graphql.ResolveParams) (interface{}, error) {
	email, _ := params.Args["email"].(string)
	password, _ := params.Args["password"].(string)
	res, err := r.ServiceFactory.UserService.Login(email, password)
	if err != nil {
		return nil, gqlerrors.FormatError(err)
	}
	jwtToken, err := jwt.GenerateToken(res)
	type authTypeResponse struct {
		Token     string
		TokenType string
	}

	return authTypeResponse{
		Token:     jwtToken,
		TokenType: "jwt",
	}, nil
}

func (r *Resolver) UpdateUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	name, _ := params.Args["name"].(string)
	email, _ := params.Args["email"].(string)
	res, err := r.ServiceFactory.UserService.Update(&entity.User{
		Id:    id,
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, gqlerrors.FormatError(err)
	}
	return res, nil
}

func (r *Resolver) Me(params graphql.ResolveParams) (interface{}, error) {
	user := params.Context.Value("currentUser").(*entity.User)
	return user, nil
}

func (r *Resolver) GetUser(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(int)
	selectedFields, err := GetSelectedFields(params)
	user, err := r.ServiceFactory.UserService.View(id, selectedFields)
	if err != nil {
		return nil, gqlerrors.FormatError(err)
	}
	return user, nil
}
func (r *Resolver) GetUsers(params graphql.ResolveParams) (interface{}, error) {
	name := params.Args["name"].(string)
	users, err := r.ServiceFactory.UserService.List(name)
	if err != nil {
		return nil, gqlerrors.FormatError(err)
	}
	return users, nil
}

func GetSelectedFields(params graphql.ResolveParams) ([]string, error) {
	fieldASTs := params.Info.FieldASTs
	if len(fieldASTs) == 0 {
		err := fmt.Errorf("GetSelectedFields: ResolveParams has no fields")
		// logrus.Error(err)
		return nil, err
	}

	fields, err := selectedFieldsFromSelections(params, fieldASTs[0].SelectionSet.Selections, ".")
	if err == nil {
		for i, f := range fields {
			fields[i] = f[1:]
		}
	}

	return fields, err
}

func selectedFieldsFromSelections(params graphql.ResolveParams, selections []ast.Selection, suffix string) ([]string, error) {
	var selected []string
	for _, s := range selections {
		switch t := s.(type) {
		case *ast.Field:
			if s.GetSelectionSet() != nil {
				var subfields []ast.Selection

				subfields = append(subfields, s.GetSelectionSet().Selections...)

				sel, err := selectedFieldsFromSelections(params, subfields, suffix+t.Name.Value+suffix)
				if err != nil {
					return nil, err
				}
				selected = append(selected, sel...)
			} else {
				selected = append(selected, suffix+s.(*ast.Field).Name.Value)
			}
		case *ast.FragmentSpread:
			n := s.(*ast.FragmentSpread).Name.Value
			frag, ok := params.Info.Fragments[n]
			if !ok {
				err := fmt.Errorf("GetSelectedFields: ResolveParams has no fields")
				// logrus.Error(err)
				return nil, err
			}
			sel, err := selectedFieldsFromSelections(params, frag.GetSelectionSet().Selections, suffix)
			if err != nil {
				return nil, err
			}
			selected = append(selected, sel...)
		default:
			err := fmt.Errorf("getSelectedFields: found unexpected selection type %v", t)
			// logrus.Error(err)
			return nil, err
		}
	}
	return selected, nil
}
