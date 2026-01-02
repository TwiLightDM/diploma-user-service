package user_service

import (
	"context"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
)

type UserHandler struct {
	userservicepb.UnimplementedUserServiceServer
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Login(_ context.Context, req *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	access, refresh, err := h.service.Login(req.Email, req.Password)
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

func (h *UserHandler) SignUp(_ context.Context, req *userservicepb.SignUpRequest) (*userservicepb.SignUpResponse, error) {
	err := h.service.Signup(req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		return &userservicepb.SignUpResponse{
			Error: err.Error(),
		}, nil
	}

	return &userservicepb.SignUpResponse{}, nil
}
