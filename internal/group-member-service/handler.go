package group_member_service

import (
	"context"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/proto/groupmemberservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GroupMemberService interface {
	CreateGroupMember(ctx context.Context, email, groupId string) (*entities.GroupMember, error)
	ReadAllGroupMembersByUserId(ctx context.Context, userId string) ([]entities.GroupMember, error)
	ReadAllGroupMembersByGroupId(ctx context.Context, groupId string) ([]entities.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id string) error
}

type GroupMemberHandler struct {
	groupmemberservicepb.UnimplementedGroupMemberServiceServer
	service GroupMemberService
}

func NewGroupMemberHandler(service GroupMemberService) *GroupMemberHandler {
	return &GroupMemberHandler{service: service}
}

func (h *GroupMemberHandler) CreateGroupMember(ctx context.Context, req *groupmemberservicepb.CreateGroupMemberRequest) (*groupmemberservicepb.CreateGroupMemberResponse, error) {
	groupMember, err := h.service.CreateGroupMember(ctx, req.Email, req.GroupId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupmemberservicepb.CreateGroupMemberResponse{
		GroupMember: &groupmemberservicepb.GroupMember{
			Id:      groupMember.Id,
			UserId:  groupMember.UserId,
			GroupId: groupMember.GroupId,
		},
	}, nil
}

func (h *GroupMemberHandler) ReadAllGroupMembersByUserId(ctx context.Context, req *groupmemberservicepb.ReadAllGroupMembersByUserIdRequest) (*groupmemberservicepb.ReadAllGroupMembersByUserIdResponse, error) {
	groupMembers, err := h.service.ReadAllGroupMembersByUserId(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	groupMembersPb := h.groupMembersToPb(groupMembers)

	return &groupmemberservicepb.ReadAllGroupMembersByUserIdResponse{
		GroupMembers: groupMembersPb,
	}, nil
}

func (h *GroupMemberHandler) ReadAllGroupMembersByGroupId(ctx context.Context, req *groupmemberservicepb.ReadAllGroupMembersByGroupIdRequest) (*groupmemberservicepb.ReadAllGroupMembersByGroupIdResponse, error) {
	groupMembers, err := h.service.ReadAllGroupMembersByGroupId(ctx, req.GroupId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	groupMembersPb := h.groupMembersToPb(groupMembers)

	return &groupmemberservicepb.ReadAllGroupMembersByGroupIdResponse{
		GroupMembers: groupMembersPb,
	}, nil
}

func (h *GroupMemberHandler) DeleteGroupMember(ctx context.Context, req *groupmemberservicepb.DeleteGroupMemberRequest) (*groupmemberservicepb.DeleteGroupMemberResponse, error) {
	err := h.service.DeleteGroupMember(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &groupmemberservicepb.DeleteGroupMemberResponse{}, nil
}

func (h *GroupMemberHandler) groupMembersToPb(groupMembers []entities.GroupMember) []*groupmemberservicepb.GroupMember {
	groupMembersPb := make([]*groupmemberservicepb.GroupMember, 0, len(groupMembers))
	for _, groupMember := range groupMembers {
		groupMembersPb = append(groupMembersPb, &groupmemberservicepb.GroupMember{
			Id:      groupMember.Id,
			UserId:  groupMember.UserId,
			GroupId: groupMember.GroupId,
		})
	}

	return groupMembersPb
}
