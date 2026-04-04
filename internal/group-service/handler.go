package group_service

import (
	"context"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/proto/groupservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupService interface {
	CreateGroup(ctx context.Context, group *entities.Group) (*entities.Group, error)
	ReedGroupById(ctx context.Context, id string) (*entities.Group, error)
	ReadAllGroupsByOwnerId(ctx context.Context, ownerId string) ([]entities.Group, error)
	UpdateGroup(ctx context.Context, group *entities.Group) (*entities.Group, error)
	DeleteGroup(ctx context.Context, id string) error
}

type GroupHandler struct {
	groupservicepb.UnimplementedGroupServiceServer
	service GroupService
}

func NewGroupHandler(service GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) CreateGroup(ctx context.Context, req *groupservicepb.CreateGroupRequest) (*groupservicepb.CreateGroupResponse, error) {
	updatedGroup, err := h.service.CreateGroup(ctx, &entities.Group{
		Title:       req.Title,
		Description: req.Description,
		OwnerId:     req.OwnerId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupservicepb.CreateGroupResponse{
		Group: &groupservicepb.Group{
			Id:          updatedGroup.Id,
			Title:       updatedGroup.Title,
			Description: updatedGroup.Description,
			OwnerId:     updatedGroup.OwnerId,
		},
	}, nil
}

func (h *GroupHandler) ReadGroup(ctx context.Context, req *groupservicepb.ReadGroupRequest) (*groupservicepb.ReadGroupResponse, error) {
	group, err := h.service.ReedGroupById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupservicepb.ReadGroupResponse{
		Group: &groupservicepb.Group{
			Id:          group.Id,
			Title:       group.Title,
			Description: group.Description,
			OwnerId:     group.OwnerId,
		},
	}, nil
}

func (h *GroupHandler) ReadAllGroupsByOwnerId(ctx context.Context, req *groupservicepb.ReadAllGroupsByOwnerIdRequest) (*groupservicepb.ReadAllGroupsByOwnerIdResponse, error) {
	groups, err := h.service.ReadAllGroupsByOwnerId(ctx, req.OwnerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	groupsPb := make([]*groupservicepb.Group, 0, len(groups))
	for _, group := range groups {
		groupsPb = append(groupsPb, &groupservicepb.Group{
			Id:          group.Id,
			Title:       group.Title,
			Description: group.Description,
			OwnerId:     group.OwnerId,
		})
	}

	return &groupservicepb.ReadAllGroupsByOwnerIdResponse{
		Groups: groupsPb,
	}, nil
}

func (h *GroupHandler) UpdateGroup(ctx context.Context, req *groupservicepb.UpdateGroupRequest) (*groupservicepb.UpdateGroupResponse, error) {
	updatedGroup, err := h.service.UpdateGroup(ctx, &entities.Group{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupservicepb.UpdateGroupResponse{
		Group: &groupservicepb.Group{
			Id:          updatedGroup.Id,
			Title:       updatedGroup.Title,
			Description: updatedGroup.Description,
			OwnerId:     updatedGroup.OwnerId,
		},
	}, nil
}

func (h *GroupHandler) DeleteGroup(ctx context.Context, req *groupservicepb.DeleteGroupRequest) (*groupservicepb.DeleteGroupResponse, error) {
	err := h.service.DeleteGroup(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupservicepb.DeleteGroupResponse{}, nil
}
