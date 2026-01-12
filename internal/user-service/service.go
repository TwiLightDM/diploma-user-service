package user_service

import (
	"context"
	"errors"
	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/TwiLightDM/diploma-user-service/package/validation"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	SignUp(ctx context.Context, email, password, fullName, role string) error
	ReedById(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdatePassword(ctx context.Context, user *entities.User) error
}

type userService struct {
	repo    UserRepository
	jwt     *utils.JWTService
	encrypt *utils.EncryptService
}

func NewUserService(repo UserRepository, jwt *utils.JWTService, encrypt *utils.EncryptService) UserService {
	return &userService{repo: repo, jwt: jwt, encrypt: encrypt}
}

func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := s.repo.ReadByEmail(ctx, email)
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

func (s *userService) SignUp(ctx context.Context, email, password, fullName, role string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	existing, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user already exists")
	}

	ok := validation.IsValidEmail(email)
	if !ok {
		return errors.New("invalid email")
	}

	ok = validation.IsStrongPassword(password)
	if !ok {
		return errors.New("invalid password")
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

	return s.repo.Create(ctx, &user)
}

func (s *userService) ReedById(ctx context.Context, id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) UpdatePassword(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error

	ok := validation.IsStrongPassword(user.Password)
	if !ok {
		return errors.New("invalid password")
	}

	user.Password, user.Salt, err = s.encrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
