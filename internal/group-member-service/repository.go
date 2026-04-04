package group_member_service

import (
	"context"
	"errors"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"gorm.io/gorm"
)

type groupMemberRepository struct {
	db *gorm.DB
}

func NewGroupMemberRepository(db *gorm.DB) GroupMemberRepository {
	return &groupMemberRepository{db: db}
}

func (r *groupMemberRepository) Create(ctx context.Context, groupMember *entities.GroupMember) error {
	return r.db.WithContext(ctx).Create(groupMember).Error
}

func (r *groupMemberRepository) ReadAllByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error) {
	var groupMembers []entities.GroupMember
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Find(&groupMembers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return groupMembers, nil
}

func (r *groupMemberRepository) ReadAllByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error) {
	var groupMembers []entities.GroupMember
	if err := r.db.
		WithContext(ctx).
		Where("group_id = ?", groupId).
		Find(&groupMembers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return groupMembers, nil
}

func (r *groupMemberRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.GroupMember{}).Error
}
