package dto

import (
	"api/db/models"
	"api/graph/gmodel"
	"time"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func MapUserToResponse(user *models.User) *gmodel.User {
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
