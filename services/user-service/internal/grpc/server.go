package grpcserver

import (
	"context"

	userv1 "github.com/novelhive/proto/user/v1"
	"github.com/novelhive/user-service/internal/domain"
	"github.com/novelhive/user-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
}

type UserProfileServiceServer interface {
	UpdateAvatar(context.Context, *structpb.Struct) (*structpb.Struct, error)
	SignInWithOAuth(context.Context, *structpb.Struct) (*structpb.Struct, error)
}

func NewUserServiceServer(uc *usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{uc: uc}
}

func (s *UserServiceServer) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	user, accessToken, refreshToken, err := s.uc.Register(domain.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &userv1.RegisterResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, accessToken, refreshToken, err := s.uc.Login(domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	return &userv1.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

func (s *UserServiceServer) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := s.uc.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, mapDomainError(err)
	}
	return &userv1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

func (s *UserServiceServer) UpdateAvatar(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	userID := stringsFromStruct(req, "user_id")
	avatarURL := stringsFromStruct(req, "avatar_url")
	if userID == "" || avatarURL == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and avatar_url are required")
	}

	user, err := s.uc.UpdateAvatarURL(userID, avatarURL)
	if err != nil {
		return nil, mapDomainError(err)
	}

	resp, err := structpb.NewStruct(map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"avatar_url": user.AvatarURL,
		"role":       user.Role,
		"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to encode profile")
	}
	return resp, nil
}

func (s *UserServiceServer) SignInWithOAuth(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	user, accessToken, refreshToken, err := s.uc.LoginWithOAuth(domain.OAuthLoginInput{
		Email:     stringsFromStruct(req, "email"),
		Username:  stringsFromStruct(req, "username"),
		AvatarURL: stringsFromStruct(req, "avatar_url"),
	})
	if err != nil {
		return nil, mapDomainError(err)
	}

	resp, err := structpb.NewStruct(map[string]interface{}{
		"user_id":       user.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"role":       user.Role,
			"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to encode oauth response")
	}
	return resp, nil
}

func stringsFromStruct(s *structpb.Struct, key string) string {
	if s == nil || s.Fields == nil {
		return ""
	}
	if value := s.Fields[key]; value != nil {
		return value.GetStringValue()
	}
	return ""
}

func RegisterUserProfileServiceServer(s grpc.ServiceRegistrar, srv UserProfileServiceServer) {
	s.RegisterService(&UserProfileService_ServiceDesc, srv)
}

func _UserProfileService_UpdateAvatar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(structpb.Struct)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).UpdateAvatar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.v1.UserProfileService/UpdateAvatar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).UpdateAvatar(ctx, req.(*structpb.Struct))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserProfileService_SignInWithOAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(structpb.Struct)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserProfileServiceServer).SignInWithOAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.v1.UserProfileService/SignInWithOAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserProfileServiceServer).SignInWithOAuth(ctx, req.(*structpb.Struct))
	}
	return interceptor(ctx, in, info, handler)
}

var UserProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserProfileService",
	HandlerType: (*UserProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateAvatar",
			Handler:    _UserProfileService_UpdateAvatar_Handler,
		},
		{
			MethodName: "SignInWithOAuth",
			Handler:    _UserProfileService_SignInWithOAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/user/v1/user_profile.proto",
}

func mapDomainError(err error) error {
	switch err {
	case domain.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrEmailExists, domain.ErrUsernameExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidPassword:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidToken, domain.ErrRefreshTokenInvalid:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrTokenExpired:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidInput:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
