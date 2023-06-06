package main

import (
	"context"
	"example/web-service-gin/graph"
	"example/web-service-gin/graph/resolver"
	"example/web-service-gin/internal/config"
	"example/web-service-gin/internal/constant/app_error"
	"example/web-service-gin/internal/service"
	"example/web-service-gin/pkg/audit"
	"example/web-service-gin/pkg/database"
	"example/web-service-gin/pkg/file"
	"example/web-service-gin/pkg/logger"
	"example/web-service-gin/pkg/middleware"
	"example/web-service-gin/pkg/middleware/auth"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// repository := repository.N(database.DB, logger.Logger)
	serviceFactory := service.InitServices(database.DB, logger.Logger)
	graphqlServer := initGraphQLServer(database.DB, serviceFactory)
	router := initRouter(graphqlServer)

	err := http.ListenAndServe(":"+port, router)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	if err != nil {
		panic(err)
	}
}

func init() {
	config.LoadEnvVariables()
	logger.InitializeLogger()
	// log.Logger.Info("Initialize database...")
	database.InitializeDatabase()
	// Register audit
	audit.RegisterAudit()
}

// Init graphql Server
func initGraphQLServer(tx *gorm.DB, serviceFactory *service.ServiceFactory) *handler.Server {

	c := graph.Config{Resolvers: &resolver.Resolver{
		ServiceFactory: serviceFactory,
	}}

	// create directive
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		_, ok := ctx.Value("currentUser").(*auth.AuthenticatedUser)
		if !ok {
			// block calling the next resolver
			return nil, fmt.Errorf("Access denied")
		}
		return next(ctx)
	}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role *string) (interface{}, error) {
		user, ok := ctx.Value("currentUser").(*auth.AuthenticatedUser)
		if !ok {
			return nil, fmt.Errorf("Access denied")
		}
		if user.Role != *role {
			return nil, fmt.Errorf("Unauthorized access")
		}

		// or let it pass through
		return next(ctx)
	}

	server := handler.NewDefaultServer(graph.NewExecutableSchema(c))
	server.Use(middleware.GqlTransaction{DB: tx})
	server.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		message := app_error.GetErrorMessage(err.Message)
		if message != "" {
			err.Extensions = map[string]interface{}{
				"code": err.Message,
			}
			err.Message = message
			return err
		}
		return err
	})

	return server
}

func initRouter(graphqlServer *handler.Server) *chi.Mux {
	router := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "https://localhost:4200", "http://localhost:4200"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", middleware.NewMiddleware(graphqlServer))
	router.HandleFunc("/download-file", file.DownloadFile)
	router.HandleFunc("/download-aspose-format", file.DownloadFormatFile)
	return router
}
