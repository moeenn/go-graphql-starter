package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5"
	"github.com/vektah/gqlparser/v2/ast"
	"graphql/config"
	dbmodels "graphql/db/models"
	"graphql/graph"
	"graphql/graph/resolvers"
	"log/slog"
	"net/http"
	"os"
)

func run(ctx context.Context) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config, err := config.NewConfig(logger)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	dbConn, err := pgx.Connect(ctx, config.Database.URI)
	if err != nil {
		return err
	}
	defer dbConn.Close(ctx)
	if err := dbConn.Ping(ctx); err != nil {
		return err
	}

	db := dbmodels.New(dbConn)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &resolvers.Resolver{
			Logger: logger,
			DB:     db,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle(config.Server.GraphiqlUrl, playground.Handler("GraphQL playground", config.Server.GraphqlUrl))
	http.Handle(config.Server.GraphqlUrl, srv)

	address := config.Server.Address()
	logger.Info("starting server", "address", address)
	return http.ListenAndServe(address, nil)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
