package user_service

import (
	"context"
	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
)

type UserHandler struct {
	userservicepb.UnimplementedUserServiceServer
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(ctx context.Context, req *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	access, refresh, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &userservicepb.LoginResponse{
			Error: err.Error(),
		}, nil
	}

	return &userservicepb.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (h *UserHandler) SignUp(ctx context.Context, req *userservicepb.SignUpRequest) (*userservicepb.SignUpResponse, error) {
	err := h.service.SignUp(ctx, req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		return &userservicepb.SignUpResponse{
			Error: err.Error(),
		}, nil
	}

	return &userservicepb.SignUpResponse{}, nil
}

func (h *UserHandler) ReadUser(ctx context.Context, req *userservicepb.ReadUserRequest) (*userservicepb.ReadUserResponse, error) {
	user, err := h.service.ReedById(ctx, req.Id)
	if err != nil {
		return &userservicepb.ReadUserResponse{
			Error: err.Error(),
		}, nil
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
		return &userservicepb.UpdateUserResponse{Error: err.Error()}, nil
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
		return &userservicepb.ChangePasswordResponse{Error: err.Error()}, nil
	}

	return &userservicepb.ChangePasswordResponse{}, nil
}
