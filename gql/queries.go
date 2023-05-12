package gql

import (
	"example/web-service-gin/controller"

	"github.com/graphql-go/graphql"
)

func GetRootQuery(contrs *controller.ServiceFactory) *graphql.Object {
	resolver := Resolver{
		ServiceFactory: contrs,
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type:        authenticatedUserType,
				Description: "Get the logged-in user's info",
				Resolve:     resolver.Me,
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
				Resolve:     resolver.GetUser,
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
				Resolve:     resolver.GetUsers,
			},
		},
	})
}
