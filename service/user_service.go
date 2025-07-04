package service

import (
	"api/db/models"
	"api/graph/gmodel"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

	//nolint: exhaustruct
	now := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}

	newUser := models.CreateUserParams{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  string(passwordHash),
		Role:      models.NullUserRole{UserRole: models.UserRoleUser, Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.DB.CreateUser(ctx, newUser); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil, errors.New("user already exists")
			}
		}

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

	users, err := s.DB.GetAllUsers(ctx, models.GetAllUsersParams{
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

	return &gmodel.UsersResponse{
		Users:      mapUserRowsToResponse(users),
		TotalCount: totalCount,
	}, nil
}

func (s Service) SetUserDeletedStatus(
	ctx context.Context,
	userID uuid.UUID,
	deleted bool,
) (*gmodel.MessageResponse, error) {
	//nolint: exhaustruct
	deletedAt := pgtype.Timestamp{Time: time.Now(), Valid: true}
	if !deleted {
		deletedAt.Valid = false
	}

	err := s.DB.SetUserDeletedStatus(ctx, models.SetUserDeletedStatusParams{
		ID:        userID,
		DeletedAt: deletedAt,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to set user deleted status: %w", err)
	}

	return &gmodel.MessageResponse{
		Message: "status updated successfully",
	}, nil
}
