package main

import (
	"context"
	"example/web-service-gin/config"
	"example/web-service-gin/graph"
	"example/web-service-gin/internal/auth"
	"example/web-service-gin/internal/constant/app_error"
	"example/web-service-gin/pkg/audit"
	"example/web-service-gin/pkg/database"
	"example/web-service-gin/pkg/file"
	"example/web-service-gin/pkg/logger"
	"example/web-service-gin/pkg/repository"
	"example/web-service-gin/pkg/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8080"

func init() {
	config.LoadEnvVariables()
	logger.InitializeLogger()
	// log.Logger.Info("Initialize database...")
	database.InitializeDatabase()
	// log.Logger.Info("Sync database...")
	database.SyncDatabase()
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repository := repository.NewGormRepository(database.DB, logger.Logger)
	serviceFactory := service.InitServices(repository)

	c := graph.Config{Resolvers: &graph.Resolver{
		ServiceFactory: serviceFactory,
	}}

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

	server.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		err.Extensions = map[string]interface{}{
			"code": err.Message,
		}
		err.Message = app_error.GetErrorMessage(err.Message)
		return err
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", AuthenticationMiddleware(PerformanceMiddleware(server)))
	http.HandleFunc("/download-file", file.DownloadFile)
	http.HandleFunc("/download-aspose-format", file.DownloadFormatFile)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func PerformanceMiddleware(next http.Handler) http.Handler {
	perfomanceFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI

		next.ServeHTTP(rw, r) // serve the original request

		duration := time.Since(start)
		// log request details
		if duration.Seconds() > 10 {
			logger.Logger.Info(fmt.Sprintf("%s: Log Duration %s\n", uri, duration))
		}
	}
	return http.HandlerFunc(perfomanceFn)
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	authenticationFn := func(rw http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(rw, r)
			return
		}

		//validate jwt token
		tokenStr := header
		user, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(rw, "Invalid token", http.StatusForbidden)
			return
		}
		// put it in context
		ctx := context.WithValue(r.Context(), "currentUser", user)
		r = r.WithContext(ctx)

		database.DB.Set(audit.CurrentUserDBScopeKey, user.Email)

		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(authenticationFn)
}
