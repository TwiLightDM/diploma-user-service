package user_service

import (
	"context"
	"fmt"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	userservicepb.UnimplementedUserServiceServer
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(ctx context.Context, req *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	accessToken, accessExpiresAt, refreshToken, RefreshExpiresAt, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userservicepb.LoginResponse{
		AccessToken:      accessToken,
		AccessExpiresAt:  timestamppb.New(*accessExpiresAt),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: timestamppb.New(*RefreshExpiresAt),
	}, nil
}

func (h *UserHandler) SignUp(ctx context.Context, req *userservicepb.SignUpRequest) (*userservicepb.SignUpResponse, error) {
	err := h.service.SignUp(ctx, req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		switch st.Code() {
		case codes.InvalidArgument:
			fmt.Println(st.Message())
		}
	}

	return &userservicepb.SignUpResponse{}, nil
}

func (h *UserHandler) ReadUser(ctx context.Context, req *userservicepb.ReadUserRequest) (*userservicepb.ReadUserResponse, error) {
	user, err := h.service.ReedById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userservicepb.ReadUserResponse{
		Email:    user.Email,
		FullName: user.FullName,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userservicepb.UpdateUserRequest) (*userservicepb.UpdateUserResponse, error) {
	updatedUser, err := h.service.UpdateUser(ctx, &entities.User{
		Id:       req.Id,
		Email:    req.Email,
		FullName: req.FullName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userservicepb.UpdateUserResponse{
		Email:    updatedUser.Email,
		FullName: updatedUser.FullName,
	}, nil
}

func (h *UserHandler) ChangePassword(ctx context.Context, req *userservicepb.ChangePasswordRequest) (*userservicepb.ChangePasswordResponse, error) {
	err := h.service.UpdatePassword(ctx, &entities.User{
		Id:       req.Id,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userservicepb.ChangePasswordResponse{}, nil
}
