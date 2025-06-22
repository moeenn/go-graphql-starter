package resolvers

import (
	"api/config"
	dbmodels "api/db/models"
	"log/slog"
)

type Resolver struct {
	Logger *slog.Logger
	DB     *dbmodels.Queries
	Config *config.Config
}
