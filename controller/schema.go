package controller

import (
	"errors"
	"example/web-service-gin/entity"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	ast "github.com/graphql-go/graphql/language/ast"
)

var authType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Auth",
	Fields: graphql.Fields{
		"tokenType": &graphql.Field{
			Type: graphql.String,
		},
		"token": &graphql.Field{
			Type: graphql.String,
		},
		"expiresIn": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
var roleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "role",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of role",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of role",
		},
	}})

var userViewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "user",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of user",
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of user",
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "The email of user",
		},
		"createdAt": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "The time that user was created",
		},
		"lastModifiedAt": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "The last time that user was updated",
		},
		"role": &graphql.Field{
			Type:        roleType,
			Description: "The role of user",
		},
	},
})

func getRootMutation(contrs *ControllerFactory) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"signup": &graphql.Field{
				Type:        authType, // the return type for this field
				Description: "Signup",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// username, _ := params.Args["username"].(string)
					// password, _ := params.Args["password"].(string)
					res := "test"
					return res, nil
				},
			},
			"login": &graphql.Field{
				Type:        authType, // the return type for this field
				Description: "Login",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					email, _ := params.Args["email"].(string)
					password, _ := params.Args["password"].(string)
					res, err := contrs.userController.Login(email, password)
					if err != nil {
						return nil, gqlerrors.FormatError(err)
					}
					tokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
					// Generate a jwt token
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"sub": res.Email,
						"exp": tokenExp,
					})
					// Sign and get the  complete encoded token as string using secret
					tokenString, err := token.SignedString([]byte(os.Getenv(("JWT_SECRET"))))
					if err != nil {
						return nil, gqlerrors.FormatError(errors.New("Generate token fail"))
					}
					type authTypeResponse struct {
						Token     string
						TokenType string
						expiresIn int64
					}
					return authTypeResponse{
						Token:     tokenString,
						TokenType: "jwt",
						expiresIn: tokenExp,
					}, nil
				},
			},
			"createUser": &graphql.Field{
				Type:        userViewType, // the return type for this field
				Description: "Create new user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					name, _ := params.Args["name"].(string)
					email, _ := params.Args["email"].(string)
					password, _ := params.Args["password"].(string)
					res, err := contrs.userController.Create(&entity.User{
						Id:       id,
						Name:     name,
						Email:    email,
						Password: password,
					})
					if err != nil {
						return nil, gqlerrors.FormatError(err)
					}
					return res, nil
				},
			},
			"updateUser": &graphql.Field{
				Type:        userViewType, // the return type for this field
				Description: "Update new user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					name, _ := params.Args["name"].(string)
					email, _ := params.Args["email"].(string)
					res, err := contrs.userController.Update(&entity.User{
						Id:    id,
						Name:  name,
						Email: email,
					})
					if err != nil {
						return nil, gqlerrors.FormatError(err)
					}
					return res, nil
				},
			},
		},
	})
}

func getRootQuery(contrs *ControllerFactory) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Me",
					Fields: graphql.Fields{
						"username": &graphql.Field{
							Type: graphql.String,
						},
					},
				}),
				Description: "Get the logged-in user's info",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// user := params.Context.Value(contextKey("currentUser")).(model.User)
					rootValue := params.Info.RootValue.(map[string]interface{})
					user := rootValue["currentUser"].(entity.User)
					return user.Name, nil
				},
			},
			"user": &graphql.Field{
				Type: userViewType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "Id of user",
					},
				},
				Description: "Get detail user by id",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(int)
					test2, err := GetSelectedFields(params)
					fmt.Println(test2)
					user, err := contrs.userController.View(id)
					if err != nil {
						return nil, gqlerrors.FormatError(err)
					}
					return user, nil
				},
			},
			"users": &graphql.Field{
				Type: graphql.NewList(userViewType),
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "The keyword search by name",
					},
				},
				Description: "Get list user",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name := params.Args["name"].(string)
					users, err := contrs.userController.List(name)
					if err != nil {
						return nil, gqlerrors.FormatError(err)
					}
					return users, nil
				},
			},
		},
	})
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

type ErrorResponse struct {
	Code    string
	Message string
}
