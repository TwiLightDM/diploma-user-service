package group_service

import (
	"context"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/google/uuid"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.Group) error
	ReadById(ctx context.Context, id string) (*entities.Group, error)
	ReadAllByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error)
	Update(ctx context.Context, group *entities.Group) (*entities.Group, error)
	Delete(ctx context.Context, id string) error
}

type groupService struct {
	repo GroupRepository
}

func NewGroupService(repo GroupRepository) GroupService {
	return &groupService{repo: repo}
}

func (s *groupService) CreateGroup(ctx context.Context, group *entities.Group) (*entities.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	group.Id = uuid.NewString()

	err := s.repo.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) ReedGroupById(ctx context.Context, id string) (*entities.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	group, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) ReadAllGroupsByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	groups, err := s.repo.ReadAllByOwnerId(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *groupService) UpdateGroup(ctx context.Context, group *entities.Group) (*entities.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var err error

	updatedGroup, err := s.repo.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

func (s *groupService) DeleteGroup(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
