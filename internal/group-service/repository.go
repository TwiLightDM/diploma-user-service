package group_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.Group) error
	ReadById(ctx context.Context, id string) (*entities.Group, error)
	ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error)
	Update(ctx context.Context, group *entities.Group) (*entities.Group, error)
	Delete(ctx context.Context, id string) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) Create(ctx context.Context, group *entities.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *groupRepository) ReadById(ctx context.Context, id string) (*entities.Group, error) {
	var group entities.Group
	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (r *groupRepository) ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error) {
	var groups []entities.Group
	if err := r.db.
		WithContext(ctx).
		Where("owner_id = ?", ownerId).
		Find(&groups).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

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

	return &updatedGroup, nil
}

func (r *groupRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.Group{}).Error
}
