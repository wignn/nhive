package grpcserver

import (
	"context"

	"github.com/novelhive/user-service/internal/domain"
	"github.com/novelhive/user-service/internal/usecase"
	userv1 "github.com/novelhive/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
}

func NewUserServiceServer(uc *usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{uc: uc}
}

func (s *UserServiceServer) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	user, token, err := s.uc.Register(domain.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &userv1.RegisterResponse{
		UserId: user.ID,
		Token:  token,
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, token, err := s.uc.Login(domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &userv1.LoginResponse{
		UserId: user.ID,
		Token:  token,
		Profile: &userv1.UserProfile{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarUrl: user.AvatarURL,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (s *UserServiceServer) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.UserProfile, error) {
	user, err := s.uc.GetProfile(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &userv1.UserProfile{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *UserServiceServer) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	userID, role, err := s.uc.ValidateToken(req.Token)
	if err != nil {
		return &userv1.ValidateTokenResponse{Valid: false}, nil
	}
	return &userv1.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Role:   role,
	}, nil
}

func (s *UserServiceServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	users, total, err := s.uc.ListUsers(page, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var profiles []*userv1.UserProfile
	for _, u := range users {
		profiles = append(profiles, &userv1.UserProfile{
			Id:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			AvatarUrl: u.AvatarURL,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return &userv1.ListUsersResponse{
		Users: profiles,
		Total: int32(total),
	}, nil
}

func (s *UserServiceServer) UpdateUserRole(ctx context.Context, req *userv1.UpdateUserRoleRequest) (*userv1.UpdateUserRoleResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	// Security: validate role at gRPC boundary
	if req.Role != "admin" && req.Role != "reader" {
		return nil, status.Error(codes.InvalidArgument, "role must be 'admin' or 'reader'")
	}
	if err := s.uc.UpdateUserRole(req.UserId, req.Role); err != nil {
		return nil, mapDomainError(err)
	}
	return &userv1.UpdateUserRoleResponse{Success: true}, nil
}

func mapDomainError(err error) error {
	switch err {
	case domain.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrEmailExists, domain.ErrUsernameExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidPassword:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidToken:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
