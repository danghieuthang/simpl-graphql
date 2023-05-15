package extensions

import (
	"context"
	"errors"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

type GqlTransaction struct {
	DB *gorm.DB
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationContextMutator
	graphql.ResponseInterceptor
} = GqlTransaction{}

// ExtensionName returns the extension name.
func (GqlTransaction) ExtensionName() string {
	return "GqlTransaction"
}

// Validate is called when adding an extension to the server, it allows validation against the servers schema.
func (t GqlTransaction) Validate(graphql.ExecutableSchema) error {
	if t.DB == nil {
		return errors.New("DBContext is null")
	}
	return nil
}

// MutateOperationContext serializes field resolvers during mutations.
func (GqlTransaction) MutateOperationContext(_ context.Context, oc *graphql.OperationContext) *gqlerror.Error {
	if op := oc.Operation; op != nil && op.Operation == ast.Mutation {
		previous := oc.ResolverMiddleware
		var mu sync.Mutex
		oc.ResolverMiddleware = func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
			mu.Lock()
			defer mu.Unlock()
			return previous(ctx, next)
		}
	}
	return nil
}

// InterceptResponse runs graphql mutations under a transaction.
func (t GqlTransaction) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if op := graphql.GetOperationContext(ctx).Operation; op == nil || op.Operation != ast.Mutation {
		return next(ctx)
	}
	tx := t.DB.WithContext(ctx).Begin()
	err := tx.Error
	if err != nil {
		return graphql.ErrorResponse(ctx,
			"cannot create transaction: %s", err.Error(),
		)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	rsp := next(ctx)
	if len(rsp.Errors) > 0 {
		_ = tx.Rollback()
		return &graphql.Response{
			Errors: rsp.Errors,
		}
	}
	if err := tx.Commit().Error; err != nil {
		return graphql.ErrorResponse(ctx,
			"cannot commit transaction: %s", err.Error(),
		)
	}
	return rsp
}
