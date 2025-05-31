package resolvers

import (
	dbmodels "graphql/db/models"
	"log/slog"
)

type Resolver struct {
	Logger *slog.Logger
	DB     *dbmodels.Queries
}
