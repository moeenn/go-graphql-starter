package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type ContextKey struct{}

var (
	JwtClaimsContextKey ContextKey = struct{}{}
)

type GlobalMiddlewareArgs struct {
	Logger                *slog.Logger
	Next                  http.Handler
	JwtSecret             []byte
	WhiteListedOperations []string
}

func (args *GlobalMiddlewareArgs) Validate() error {
	if args.Logger == nil {
		return errors.New("logger is required")
	}

	if args.Next == nil {
		return errors.New("next function is required")
	}

	if len(args.JwtSecret) == 0 {
		return errors.New("jwt secret cannot be an empty string")
	}

	if args.WhiteListedOperations == nil {
		args.WhiteListedOperations = []string{}
	}

	return nil
}

func NewGraphqlGlobalMiddleware(args *GlobalMiddlewareArgs) (http.Handler, error) {
	if err := args.Validate(); err != nil {
		return nil, fmt.Errorf("invalid global middleware args: %w", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operations, forward, err := parseGraphqlRequestOperations(r)
		if forward {
			args.Next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			err = fmt.Errorf("failed to detect graphql operations: %w", err)
			reportHttpError(w, err, http.StatusUnauthorized)
			return
		}

		allOpsWhitelisted := allOperationsWhitelisted(operations, args.WhiteListedOperations)
		if allOpsWhitelisted {
			args.Next.ServeHTTP(w, r)
			return
		}

		bearerToken, err := readBearerToken(r)
		if err != nil {
			reportHttpError(w, err, http.StatusBadRequest)
			return
		}

		parsedClaims, err := validateAndParseJwtClaims(args.JwtSecret, bearerToken)
		if err != nil {
			reportHttpError(w, fmt.Errorf("unauthorized: %w", err), http.StatusUnauthorized)
			return
		}

		// allow the incoming request.
		ctx := context.WithValue(r.Context(), JwtClaimsContextKey, parsedClaims)
		args.Next.ServeHTTP(w, r.WithContext(ctx))
	})

	return handler, nil
}
