package service

import (
	"api/graph/gmodel"
	"api/internal/models"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func mapUserToResponse(user *models.User) *gmodel.User {
	status := gmodel.UserStatusActive
	if user.DeletedAt.Valid {
		status = gmodel.UserStatusInactive
	}

	return &gmodel.User{
		ID:        user.Id,
		Email:     user.Email,
		Status:    status,
		CreatedAt: formatTime(user.CreatedAt),
		UpdatedAt: formatTime(user.UpdatedAt),
	}
}

func mapUserRowToResponse(row *models.User) *gmodel.User {
	status := gmodel.UserStatusActive
	if row.DeletedAt.Valid {
		status = gmodel.UserStatusInactive
	}

	return &gmodel.User{
		ID:        row.Id,
		Email:     row.Email,
		Status:    status,
		CreatedAt: formatTime(row.CreatedAt),
		UpdatedAt: formatTime(row.UpdatedAt),
	}
}

func mapUserRowsToResponse(rows []*models.User) []*gmodel.User {
	size := len(rows)
	res := make([]*gmodel.User, size)
	for i := range size {
		res[i] = mapUserRowToResponse(rows[i])
	}
	return res
}
