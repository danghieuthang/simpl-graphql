package graph

//go:generate go run github.com/99designs/gqlgen generate

import "example/web-service-gin/pkg/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ServiceFactory *service.ServiceFactory
}
