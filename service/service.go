package service

import (
	"api/config"
	dbmodels "api/db/models"
	"log/slog"
)

type Service struct {
	Logger *slog.Logger
	DB     *dbmodels.Queries
	Config *config.Config
}
