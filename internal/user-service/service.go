package user_service

import (
	"context"
	"errors"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	ReadByEmail(ctx context.Context, email string) (*entities.User, error)
	ReadById(ctx context.Context, id string) (*entities.User, error)
	ReadAll(ctx context.Context) ([]entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
}

type userService struct {
	repo     UserRepository
	validate utils.ValidationService
	jwt      utils.JWTService
	encrypt  utils.EncryptService
}

func NewUserService(repo UserRepository, validate utils.ValidationService, jwt utils.JWTService, encrypt utils.EncryptService) UserService {
	return &userService{repo: repo, validate: validate, jwt: jwt, encrypt: encrypt}
}

func (s *userService) Login(ctx context.Context, email, password string) (string, *time.Time, string, *time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return "", nil, "", nil, err
	}
	if user == nil {
		return "", nil, "", nil, errors.New("user not found")
	}

	if err = s.encrypt.PasswordComparison(user.Password, password, user.Salt); err != nil {
		return "", nil, "", nil, err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, accessExpiresAt, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return "", nil, "", nil, err
	}

	refreshToken, refreshExpiresAt, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return "", nil, "", nil, err
	}

	return accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, nil
}

func (s *userService) SignUp(ctx context.Context, user *entities.User) (*entities.User, string, *time.Time, string, *time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	existing, err := s.repo.ReadByEmail(ctx, user.Email)
	if err != nil {
		return nil, "", nil, "", nil, err
	}
	if existing != nil {
		return nil, "", nil, "", nil, errors.New("user already exists")
	}

	ok := s.validate.IsValidEmail(user.Email)
	if !ok {
		return nil, "", nil, "", nil, errors.New("invalid email")
	}

	ok = s.validate.IsStrongPassword(user.Password)
	if !ok {
		return nil, "", nil, "", nil, errors.New("invalid password")
	}

	hashedPassword, salt, err := s.encrypt.HashPassword(user.Password)
	if err != nil {
		return nil, "", nil, "", nil, err
	}

	user.Id = uuid.NewString()
	user.Password = hashedPassword
	user.Salt = salt

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", nil, "", nil, err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, accessExpiresAt, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return nil, "", nil, "", nil, err
	}

	refreshToken, refreshExpiresAt, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return nil, "", nil, "", nil, err
	}

	return user, accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, nil
}

func (s *userService) Refresh(_ context.Context, id, role string) (string, *time.Time, string, *time.Time, error) {
	data := make(map[string]any)
	data["id"] = id
	data["role"] = role

	accessToken, accessExpiresAt, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return "", nil, "", nil, err
	}

	refreshToken, refreshExpiresAt, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return "", nil, "", nil, err
	}

	return accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, nil
}

func (s *userService) ReadUserById(ctx context.Context, id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	user, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ReadAll(ctx context.Context) ([]entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	users, err := s.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) UpdatePassword(ctx context.Context, user *entities.User) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	ok := s.validate.IsStrongPassword(user.Password)
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
