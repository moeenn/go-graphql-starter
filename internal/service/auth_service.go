package service

import (
	"api/graph/gmodel"
	"api/internal/helpers/jwt"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(ctx context.Context, input gmodel.LoginInput) (*gmodel.LoginResponse, error) {
	user, err := s.DB.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, expiry, err := jwt.NewExpiringToken(&jwt.ExpiringTokenArgs{
		UserId:        user.Id,
		Email:         user.Email,
		Role:          string(user.Role),
		JwtSecret:     s.Config.Auth.JwtSecret,
		ExpiryMinutes: s.Config.Auth.JwtExpiryMinutes,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	res := &gmodel.LoginResponse{
		User: mapUserToResponse(user),
		Token: &gmodel.UserToken{
			AccessToken: token,
			Expiry:      expiry,
		},
	}

	return res, nil
}
