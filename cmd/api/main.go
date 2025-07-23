package main

import (
	"api/graph"
	"api/graph/directives"
	"api/graph/resolvers"
	"api/internal/config"
	"api/internal/middleware"
	"api/internal/persistence"
	"api/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jmoiron/sqlx"
	"github.com/vektah/gqlparser/v2/ast"
)

func run(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config, err := config.NewConfig(logger)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// connect to database.
	db, err := sqlx.Open("postgres", config.Database.URI)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close db connection",
				"error", err.Error(),
			)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	persistence := persistence.NewPersistence(db, logger)

	//nolint: exhaustruct
	graphqlConfig := graph.Config{
		Resolvers: &resolvers.Resolver{
			Service: &service.Service{
				Logger: logger,
				DB:     persistence,
				Config: config,
			},
		},
	}
	graphqlConfig.Directives.HasRole = directives.HasRoleDirective(logger)

	srv := handler.New(graph.NewExecutableSchema(graphqlConfig))

	//nolint: exhaustruct
	srv.AddTransport(transport.Options{})
	//nolint: exhaustruct
	srv.AddTransport(transport.GET{})
	//nolint: exhaustruct
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	globalMiddleware, err := middleware.NewGraphqlGlobalMiddleware(
		&middleware.GlobalMiddlewareArgs{
			Logger:                logger,
			Next:                  srv,
			JwtSecret:             config.Auth.JwtSecret,
			WhiteListedOperations: config.Auth.WhiteListedOperations,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to initialize global middleware: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(config.Server.GraphiqlUrl, playground.Handler("GraphQL playground", config.Server.GraphqlUrl))
	mux.Handle(config.Server.GraphqlUrl, globalMiddleware)

	address := config.Server.Address()
	logger.Info("starting server", "address", address)

	//nolint: exhaustruct
	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadTimeout:       config.Server.Timeout,
		WriteTimeout:      config.Server.Timeout,
		IdleTimeout:       config.Server.Timeout,
		ReadHeaderTimeout: config.Server.Timeout,
	}

	return server.ListenAndServe()
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
