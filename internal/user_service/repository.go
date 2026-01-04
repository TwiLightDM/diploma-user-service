package user_service

import (
	"context"
	"errors"
	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	ReadByEmail(ctx context.Context, email string) (*entities.User, error)
	ReadByID(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) ReadByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.
		WithContext(ctx).
		Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ReadByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.
		WithContext(ctx).
		First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	var updatedUser entities.User
	err := r.db.
		WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", user.Id).
		Updates(user).
		Scan(&updatedUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &updatedUser, nil
}
