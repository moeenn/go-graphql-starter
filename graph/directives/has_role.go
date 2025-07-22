package directives

import (
	"api/graph/gmodel"
	"api/internal/middleware"
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
)

type Directive func(ctx context.Context, obj any, next graphql.Resolver, role gmodel.Role) (res any, err error)

func HasRoleDirective(logger *slog.Logger) Directive {
	return func(ctx context.Context, obj any, next graphql.Resolver, role gmodel.Role) (any, error) {
		jwtClaim, err := middleware.GetJwtClaims(ctx)
		if err != nil {
			logger.Warn("missing jwt claims in context", "error", err)
			return nil, errors.New("failed to get jwt claims from context")
		}

		if jwtClaim.Role != string(role) {
			return nil, errors.New("you don't have access to the requested resource")
		}

		return next(ctx)
	}
}
