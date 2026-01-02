package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTService struct {
	Key             string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func NewJWTService(key string, accessDuration, refreshDuration time.Duration) *JWTService {
	return &JWTService{
		Key:             key,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	}
}

func (s *JWTService) generateJWT(data map[string]any, expiration int64) (string, error) {
	claims := jwt.MapClaims{
		"exp": expiration,
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
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *JWTService) GenerateRefreshJWT(data map[string]any) (string, error) {
	expiration := time.Now().Add(s.RefreshDuration).Unix()
	delete(data, "exp")
	return s.generateJWT(data, expiration)

}

func (s *JWTService) GenerateAccessJWT(data map[string]any) (string, error) {
	expiration := time.Now().Add(s.AccessDuration).Unix()
	delete(data, "exp")
	return s.generateJWT(data, expiration)
}
