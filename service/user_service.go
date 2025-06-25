package service

import (
	"api/db/models"
	"api/graph/gmodel"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) CreateAccount(ctx context.Context, input gmodel.CreateAccountInput) (*gmodel.CreateAccountResponse, error) {
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("password confirmation failed")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &gmodel.CreateAccountResponse{
		Message: "Account created successfully",
	}, nil
}
