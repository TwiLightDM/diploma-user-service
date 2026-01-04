package user_service

import (
	"context"
	"errors"
	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	SignUp(ctx context.Context, email, password, fullName, role string) error
	ReedById(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
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

	user, err := s.repo.ReadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error
	if user.Password != "" {
		user.Password, user.Salt, err = s.encrypt.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
	}

	updatedUser, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
