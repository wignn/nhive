package clients

import (
	"context"
	"log"

	"github.com/novelhive/pkg/grpcauth"
	commentv1 "github.com/novelhive/proto/comment/v1"
	libraryv1 "github.com/novelhive/proto/library/v1"
	notificationv1 "github.com/novelhive/proto/notification/v1"
	novelv1 "github.com/novelhive/proto/novel/v1"
	userv1 "github.com/novelhive/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type Clients struct {
	Novel        novelv1.NovelServiceClient
	User         userv1.UserServiceClient
	Comment      commentv1.CommentServiceClient
	Library      libraryv1.LibraryServiceClient
	Notification notificationv1.NotificationServiceClient
	UserConn     *grpc.ClientConn

	conns []*grpc.ClientConn
}

func New(userAddr, novelAddr, commentAddr, libraryAddr, notificationAddr, apiKey string) *Clients {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(grpcauth.NewCredentials(apiKey)),
	}

	userConn := mustDial(userAddr, opts)
	novelConn := mustDial(novelAddr, opts)
	commentConn := mustDial(commentAddr, opts)
	libraryConn := mustDial(libraryAddr, opts)
	notificationConn := mustDial(notificationAddr, opts)

	return &Clients{
		User:         userv1.NewUserServiceClient(userConn),
		Novel:        novelv1.NewNovelServiceClient(novelConn),
		Comment:      commentv1.NewCommentServiceClient(commentConn),
		Library:      libraryv1.NewLibraryServiceClient(libraryConn),
		Notification: notificationv1.NewNotificationServiceClient(notificationConn),
		UserConn:     userConn,
		conns:        []*grpc.ClientConn{userConn, novelConn, commentConn, libraryConn, notificationConn},
	}
}

func (c *Clients) UpdateUserAvatar(ctx context.Context, userID, avatarURL string) (map[string]interface{}, error) {
	req, err := structpb.NewStruct(map[string]interface{}{
		"user_id":    userID,
		"avatar_url": avatarURL,
	})
	if err != nil {
		return nil, err
	}

	resp := new(structpb.Struct)
	if err := c.UserConn.Invoke(ctx, "/user.v1.UserProfileService/UpdateAvatar", req, resp); err != nil {
		return nil, err
	}
	return resp.AsMap(), nil
}

func (c *Clients) SignInWithOAuth(ctx context.Context, email, username, avatarURL string) (map[string]interface{}, error) {
	req, err := structpb.NewStruct(map[string]interface{}{
		"email":      email,
		"username":   username,
		"avatar_url": avatarURL,
	})
	if err != nil {
		return nil, err
	}

	resp := new(structpb.Struct)
	if err := c.UserConn.Invoke(ctx, "/user.v1.UserProfileService/SignInWithOAuth", req, resp); err != nil {
		return nil, err
	}
	return resp.AsMap(), nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		conn.Close()
	}
}

func mustDial(addr string, opts []grpc.DialOption) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", addr, err)
	}
	return conn
}
