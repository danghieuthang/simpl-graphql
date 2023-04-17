package gql

import "github.com/graphql-go/graphql"

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

var authenticatedUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Me",
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})
