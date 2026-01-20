package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateRefreshJWT(data map[string]any) (string, *time.Time, error)
	GenerateAccessJWT(data map[string]any) (string, *time.Time, error)
}

type jwtService struct {
	Key             string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func NewJWTService(key string, accessDuration, refreshDuration time.Duration) JWTService {
	return &jwtService{
		Key:             key,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	}
}

func (s *jwtService) generateJWT(data map[string]any, expiresAt *time.Time) (string, *time.Time, error) {
	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
	}

	for key, value := range data {
		if key == "id" {
			claims["sub"] = value
		} else {
			claims[key] = value
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.Key))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, expiresAt, nil
}

func (s *jwtService) GenerateRefreshJWT(data map[string]any) (string, *time.Time, error) {
	expiresAt := time.Now().Add(s.RefreshDuration)
	delete(data, "exp")

	return s.generateJWT(data, &expiresAt)
}

func (s *jwtService) GenerateAccessJWT(data map[string]any) (string, *time.Time, error) {
	expiresAt := time.Now().Add(s.AccessDuration)
	delete(data, "exp")

	return s.generateJWT(data, &expiresAt)
}
