package group_service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type groupRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGroupRepository(db *gorm.DB, redis *redis.Client) GroupRepository {
	return &groupRepository{
		db:    db,
		redis: redis,
	}
}

const groupCacheTTL = 10 * time.Minute

func (r *groupRepository) Create(ctx context.Context, group *entities.Group) error {
	err := r.db.WithContext(ctx).Create(group).Error
	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx, "groups:owner:"+group.OwnerId).Err()

	return nil
}

func (r *groupRepository) ReadById(ctx context.Context, id string) (*entities.Group, error) {
	cacheKey := "group:" + id

	cachedGroup, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var group entities.Group

		if json.Unmarshal([]byte(cachedGroup), &group) == nil {
			return &group, nil
		}
	}

	var group entities.Group

	if err = r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&group).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(group)
	_ = r.redis.Set(ctx, cacheKey, data, groupCacheTTL).Err()

	return &group, nil
}

func (r *groupRepository) ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error) {
	cacheKey := "groups:owner:" + ownerId

	cachedGroups, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var groups []entities.Group

		if json.Unmarshal([]byte(cachedGroups), &groups) == nil {
			return groups, nil
		}
	}

	var groups []entities.Group

	if err = r.db.
		WithContext(ctx).
		Where("owner_id = ?", ownerId).
		Find(&groups).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	data, _ := json.Marshal(groups)
	_ = r.redis.Set(ctx, cacheKey, data, groupCacheTTL).Err()

	return groups, nil
}

func (r *groupRepository) Update(ctx context.Context, group *entities.Group) (*entities.Group, error) {
	var updatedGroup entities.Group

	err := r.db.
		WithContext(ctx).
		Model(&entities.Group{}).
		Where("id = ?", group.Id).
		Updates(group).
		Scan(&updatedGroup).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	_ = r.redis.Del(ctx,
		"group:"+group.Id,
		"groups:owner:"+group.OwnerId,
	).Err()

	return &updatedGroup, nil
}

func (r *groupRepository) Delete(ctx context.Context, id string) error {
	group, err := r.ReadById(ctx, id)
	if err != nil {
		return err
	}

	if group == nil {
		return nil
	}

	err = r.db.
		WithContext(ctx).
		Where("id = ?", id).
		Delete(&entities.Group{}).Error

	if err != nil {
		return err
	}

	_ = r.redis.Del(ctx,
		"group:"+id,
		"groups:owner:"+group.OwnerId,
	).Err()

	return nil
}
