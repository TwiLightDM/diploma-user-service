package user_service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{
		db:    db,
		redis: redis,
	}
}

const userCacheTTL = 10 * time.Minute

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx, "users:all").Err()

	return nil
}

func (r *userRepository) ReadByEmail(ctx context.Context, email string) (*entities.User, error) {
	cacheKey := "user:email:" + email

	cachedUser, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user entities.User

		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil
		}
	}

	var user entities.User

	if err = r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(user)
	_ = r.redis.Set(ctx, cacheKey, data, userCacheTTL).Err()

	return &user, nil
}

func (r *userRepository) ReadById(ctx context.Context, id string) (*entities.User, error) {
	cacheKey := "user:" + id

	cachedUser, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user entities.User

		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil
		}
	}

	var user entities.User

	if err = r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(user)
	_ = r.redis.Set(ctx, cacheKey, data, userCacheTTL).Err()

	return &user, nil
}

func (r *userRepository) ReadAll(ctx context.Context) ([]entities.User, error) {
	const cacheKey = "users:all"

	cachedUsers, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []entities.User

		if json.Unmarshal([]byte(cachedUsers), &users) == nil {
			return users, nil
		}
	}

	var users []entities.User

	if err = r.db.
		WithContext(ctx).
		Find(&users).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(users)
	_ = r.redis.Set(ctx, cacheKey, data, userCacheTTL).Err()

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
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

	_ = r.redis.Del(ctx,
		"user:"+user.Id,
		"user:email:"+user.Email,
		"users:all",
	).Err()

	return &updatedUser, nil
}
