package utils

import (
	"regexp"
	"unicode"
)

const MinPasswordSize = 8

var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+-]+@([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`,
)

type ValidationService interface {
	IsValidEmail(email string) bool
	IsStrongPassword(password string) bool
}

type validationService struct {
}

func NewValidationService() ValidationService {
	return &validationService{}
}

func (v *validationService) IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (v *validationService) IsStrongPassword(password string) bool {
	if len(password) < MinPasswordSize {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
