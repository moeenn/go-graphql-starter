package service

import (
	"api/db/models"
	"api/graph/gmodel"
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
		ID:        user.ID.String(),
		Email:     user.Email,
		Status:    status,
		CreatedAt: formatTime(user.CreatedAt.Time),
		UpdatedAt: formatTime(user.UpdatedAt.Time),
	}
}

func mapUserRowToResponse(row models.GetAllUsersRow) *gmodel.User {
	status := gmodel.UserStatusActive
	if row.DeletedAt.Valid {
		status = gmodel.UserStatusInactive
	}

	return &gmodel.User{
		ID:        row.ID.String(),
		Email:     row.Email,
		Status:    status,
		CreatedAt: formatTime(row.CreatedAt.Time),
		UpdatedAt: formatTime(row.UpdatedAt.Time),
	}
}

func mapUserRowsToResponse(rows []models.GetAllUsersRow) []*gmodel.User {
	size := len(rows)
	res := make([]*gmodel.User, size)
	for i := range size {
		res[i] = mapUserRowToResponse(rows[i])
	}
	return res
}
