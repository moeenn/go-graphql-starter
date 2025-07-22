package persistence

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Persistence struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewPersistence(db *sqlx.DB, logger *slog.Logger) *Persistence {
	return &Persistence{db, logger}
}
