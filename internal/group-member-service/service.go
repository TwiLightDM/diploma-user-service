package group_member_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/google/uuid"
)

type GroupMemberService interface {
	CreateGroupMember(ctx context.Context, groupMember *entities.GroupMember) (*entities.GroupMember, error)
	ReadAllGroupMembersByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error)
	ReadAllGroupMembersByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id string) error
}

type groupMemberService struct {
	repo GroupMemberRepository
}

func NewGroupMemberService(repo GroupMemberRepository) GroupMemberService {
	return &groupMemberService{repo: repo}
}

func (s *groupMemberService) CreateGroupMember(ctx context.Context, groupMember *entities.GroupMember) (*entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupMember.Id = uuid.NewString()

	err := s.repo.Create(ctx, groupMember)
	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (s *groupMemberService) ReadAllGroupMembersByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupMember, err := s.repo.ReadAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (s *groupMemberService) ReadAllGroupMembersByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	groupMembers, err := s.repo.ReadAllByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	return groupMembers, nil
}

func (s *groupMemberService) DeleteGroupMember(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
