package service

import (
	"api/config"
	dbmodels "api/db/models"
	"fmt"
	"log/slog"
	"math"
)

type Service struct {
	Logger *slog.Logger
	DB     *dbmodels.Queries
	Config *config.Config
}

func safeInt64ToInt32(val int64) (int32, error) {
	if val > math.MaxInt32 || val < math.MinInt32 {
		return 0, fmt.Errorf("value %d out of int32 range %d", val, math.MaxInt32)
	}
	return int32(val), nil
}

func ParseLimitOffset(limit, offset int64) (int32, int32, error) {
	parsedLimit, err := safeInt64ToInt32(limit)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse limit: %w", err)
	}

	parsedOffset, err := safeInt64ToInt32(offset)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse offset: %w", err)
	}

	return parsedLimit, parsedOffset, nil
}
