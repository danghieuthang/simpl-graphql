package middleware

import (
	"context"
	"example/web-service-gin/internal/auth"
	"example/web-service-gin/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Define the performance middleware
func performanceMiddleware(next http.Handler) http.Handler {
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

// Define the authentication middleware
func authenticationMiddleware(next http.Handler) http.Handler {
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
		// create a new request context contain the authenticated user
		ctx := context.WithValue(r.Context(), "currentUser", user)

		// Create a new request using new context
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
	return http.HandlerFunc(authenticationFn)
}

func NewMiddleware(server *handler.Server, handlers ...interface{}) http.Handler {
	return performanceMiddleware(authenticationMiddleware(server))
}
