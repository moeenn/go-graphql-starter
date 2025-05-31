package main

import (
	"fmt"
	"graphql/config"
	"graphql/graph"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"graphql/graph/resolvers"

	"database/sql"
	dbmodels "graphql/db/models"
	_ "modernc.org/sqlite"
)

func run() error {
	config, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	dbConn, err := sql.Open("sqlite", config.Database.FilePath)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	db := dbmodels.New(dbConn)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	address := config.Server.Address()
	logger.Info("starting server", "address", address)
	return http.ListenAndServe(address, nil)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
