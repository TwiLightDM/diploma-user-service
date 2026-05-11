package group_member_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/google/uuid"
)

type GroupMemberRepository interface {
	Create(ctx context.Context, groupMember *entities.GroupMember) error
	ReadAllByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error)
	ReadAllByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	ReadByEmail(ctx context.Context, email string) (*entities.User, error)
}

type groupMemberService struct {
	repo     GroupMemberRepository
	userRepo UserRepository
}

func NewGroupMemberService(repo GroupMemberRepository, userRepo UserRepository) GroupMemberService {
	return &groupMemberService{repo: repo, userRepo: userRepo}
}

func (s *groupMemberService) CreateGroupMember(ctx context.Context, email, groupId string) (*entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	user, err := s.userRepo.ReadByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	groupMember := &entities.GroupMember{
		GroupId: groupId,
		UserId:  user.Id,
	}
	groupMember.Id = uuid.NewString()

	err = s.repo.Create(ctx, groupMember)
	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (s *groupMemberService) ReadAllGroupMembersByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	groupMember, err := s.repo.ReadAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (s *groupMemberService) ReadAllGroupMembersByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	groupMembers, err := s.repo.ReadAllByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	return groupMembers, nil
}

func (s *groupMemberService) DeleteGroupMember(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
