package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type EncryptService struct {
	SaltLength int
}

func NewEncryptionService(salt int) *EncryptService {
	return &EncryptService{
		SaltLength: salt,
	}
}

func (e EncryptService) HashPassword(password string) (string, string, error) {
	salt, err := e.saltGeneration()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.MinCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), salt, nil
}

func (e EncryptService) PasswordComparison(hashedPassword, password, salt string) error {
	saltPassword := salt + password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltPassword))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	return nil
}

func (e EncryptService) saltGeneration() (string, error) {
	bytes := make([]byte, e.SaltLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes)[:e.SaltLength], nil
}
