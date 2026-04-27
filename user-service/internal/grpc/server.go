package grpcserver

import (
	"context"

	"github.com/novelhive/user-service/internal/domain"
	"github.com/novelhive/user-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserServiceServer implements the gRPC UserService.
// Proto-generated interface is defined inline here since we're not using protoc-gen yet.
type UserServiceServer struct {
	uc *usecase.UserUsecase
	UnimplementedUserServiceServer
}

// We define the gRPC server interface manually to avoid proto-gen dependency during development.
// In production, replace with generated code from proto/user/v1/user.proto

type UnimplementedUserServiceServer struct{}

func NewUserServiceServer(uc *usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{uc: uc}
}

// RegisterRequest/Response types aligned with proto definitions
type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type RegisterResponse struct {
	UserId string
	Token  string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserId  string
	Token   string
	Profile *UserProfile
}

type GetProfileRequest struct {
	UserId string
}

type ValidateTokenRequest struct {
	Token string
}

type ValidateTokenResponse struct {
	Valid  bool
	UserId string
	Role   string
}

type UserProfile struct {
	Id        string
	Username  string
	Email     string
	AvatarUrl string
	Role      string
	CreatedAt string
}

func (s *UserServiceServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	user, token, err := s.uc.Register(domain.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &RegisterResponse{
		UserId: user.ID,
		Token:  token,
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, token, err := s.uc.Login(domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &LoginResponse{
		UserId: user.ID,
		Token:  token,
		Profile: &UserProfile{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarUrl: user.AvatarURL,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func (s *UserServiceServer) GetProfile(ctx context.Context, req *GetProfileRequest) (*UserProfile, error) {
	user, err := s.uc.GetProfile(req.UserId)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &UserProfile{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarUrl: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *UserServiceServer) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	userID, role, err := s.uc.ValidateToken(req.Token)
	if err != nil {
		return &ValidateTokenResponse{Valid: false}, nil
	}
	return &ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Role:   role,
	}, nil
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

// Ensure unused import is used
var _ = emptypb.Empty{}
