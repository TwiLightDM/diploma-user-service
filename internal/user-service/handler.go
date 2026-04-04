package user_service

import (
	"context"
	"fmt"
	"time"

	"github.com/TwiLightDM/diploma-user-service/internal/entities"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, *time.Time, string, *time.Time, error)
	SignUp(ctx context.Context, user *entities.User) (*entities.User, string, *time.Time, string, *time.Time, error)
	ReedUserById(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdatePassword(ctx context.Context, user *entities.User) error
}

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
	user, accessToken, accessExpiresAt, refreshToken, RefreshExpiresAt, err := h.service.SignUp(ctx, &entities.User{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     req.Role,
	})

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

	return &userservicepb.SignUpResponse{
		AccessToken:      accessToken,
		AccessExpiresAt:  timestamppb.New(*accessExpiresAt),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: timestamppb.New(*RefreshExpiresAt),
		User: &userservicepb.User{
			Id:       user.Id,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}

func (h *UserHandler) ReadUser(ctx context.Context, req *userservicepb.ReadUserRequest) (*userservicepb.ReadUserResponse, error) {
	user, err := h.service.ReedUserById(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userservicepb.ReadUserResponse{
		User: &userservicepb.User{
			Id:       user.Id,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
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
		User: &userservicepb.User{
			Id:       updatedUser.Id,
			Email:    updatedUser.Email,
			FullName: updatedUser.FullName,
			Role:     updatedUser.Role,
		},
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
