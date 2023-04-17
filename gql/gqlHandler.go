package gql

import (
	"context"
	"example/web-service-gin/controller"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Schema builds a graphql schema and returns it
func Schema(controllers *controller.ControllerFactory) graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    GetRootQuery(controllers),
		Mutation: GetRootMutation(controllers),
	})
	if err != nil {
		panic(err)
	}

	return schema
}

// GraphqlHandlfunc is a handler for the graphql endpoint.
func GraphqlHandlfunc(schema graphql.Schema) *handler.Handler {
	return handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
		RootObjectFn: func(ctx context.Context, req *http.Request) map[string]interface{} {
			// token := req.Header.Get("token")
			// user := entity.User{
			// 	Id:   123,
			// 	Name: "Init user",
			// }
			return map[string]interface{}{
				// "currentUser": user,
			}
		},
	})
}
