package service

import (
	"api/graph/gmodel"
	"api/internal/models"
	"api/internal/persistence"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) CreateAccount(
	ctx context.Context,
	input gmodel.CreateAccountInput,
) (*gmodel.MessageResponse, error) {
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("password confirmation failed")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	newUser := &models.User{
		Id:        uuid.NewString(),
		Email:     input.Email,
		Password:  string(passwordHash),
		Role:      models.UserRoleUser,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: sql.NullTime{Valid: false, Time: time.Now()},
	}

	if err := s.DB.CreateUser(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &gmodel.MessageResponse{
		Message: "account created successfully",
	}, nil
}

func (s Service) GetAllUsers(
	ctx context.Context,
	limit int64,
	offset int64,
) (*gmodel.UsersResponse, error) {

	parsedLimit, parsedOffset, err := ParseLimitOffset(limit, offset)
	if err != nil {
		return nil, err
	}

	users, err := s.DB.ListAllUsers(ctx, &persistence.ListAllUsersArgs{
		Limit:  parsedLimit,
		Offset: parsedOffset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var totalCount int64 = 0
	if len(users) > 0 {
		totalCount = users[0].TotalCount
	}

	userModels := make([]*models.User, len(users))
	for i := range len(users) {
		userModels[i] = &users[i].User
	}

	return &gmodel.UsersResponse{
		Users:      mapUserRowsToResponse(userModels),
		TotalCount: totalCount,
	}, nil
}

func (s Service) SetUserDeletedStatus(
	ctx context.Context,
	userID uuid.UUID,
	deleted bool,
) (*gmodel.MessageResponse, error) {
	//nolint: exhaustruct
	deletedAt := sql.NullTime{Time: time.Now(), Valid: true}
	if !deleted {
		deletedAt.Valid = false
	}

	args := &persistence.SetUserDeleteStatusArgs{
		UserId:    userID.String(),
		DeletedAt: deletedAt,
	}

	err := s.DB.SetUserDeleteStatus(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to set user deleted status: %w", err)
	}

	return &gmodel.MessageResponse{
		Message: "status updated successfully",
	}, nil
}
