package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId string
	Email  string
}

func validateAndParseJwtClaims(jwtSecret []byte, bearerToken string) (*JwtClaims, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}

	// Check token expiration
	exp, err := claims.GetExpirationTime()
	if err != nil || exp.Time.Before(time.Now()) {
		return nil, errors.New("expired jwt")
	}

	// Access claims
	parsedClaims, err := jwtClaimsFromMap(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt claims: %w", err)
	}

	return parsedClaims, nil
}

func jwtClaimsFromMap(claims map[string]any) (*JwtClaims, error) {
	userId, userIdOk := claims["userId"].(string)
	if !userIdOk {
		return nil, errors.New("invalid or missing userId")
	}

	email, emailOk := claims["email"].(string)
	if !emailOk {
		return nil, errors.New("invalid or missing email")
	}

	return &JwtClaims{
		UserId: userId,
		Email:  email,
	}, nil
}

func readBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return authHeader, errors.New("missing authorization header")
	}

	bearerErr := errors.New("authorization header does not contain a valid bearer token")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", bearerErr
	}

	pieces := strings.Split(authHeader, " ")
	if len(pieces) != 2 {
		return "", bearerErr
	}

	token := pieces[1]
	return token, nil
}
