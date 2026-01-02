package user_service

import (
	"errors"
	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/google/uuid"
)

type UserService interface {
	Login(email, password string) (string, string, error)
	Signup(email, password, fullName, role string) error
}

type userService struct {
	repo    UserRepository
	jwt     *utils.JWTService
	encrypt *utils.EncryptService
}

func NewUserService(repo UserRepository, jwt *utils.JWTService, encrypt *utils.EncryptService) UserService {
	return &userService{repo: repo, jwt: jwt, encrypt: encrypt}
}

func (s *userService) Login(email, password string) (string, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	if err = s.encrypt.PasswordComparison(user.Password, password, user.Salt); err != nil {
		return "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) Signup(email, password, fullName, role string) error {
	existing, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user already exists")
	}

	hashedPassword, salt, err := s.encrypt.HashPassword(password)
	if err != nil {
		return err
	}

	user := entities.User{
		Id:       uuid.NewString(),
		FullName: fullName,
		Email:    email,
		Password: hashedPassword,
		Salt:     salt,
		Role:     role,
	}

	return s.repo.Create(&user)
}
