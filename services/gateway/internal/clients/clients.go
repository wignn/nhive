package clients

import (
	"log"

	commentv1 "github.com/novelhive/proto/comment/v1"
	libraryv1 "github.com/novelhive/proto/library/v1"
	novelv1 "github.com/novelhive/proto/novel/v1"
	commentv1 "github.com/novelhive/proto/comment/v1"
	libraryv1 "github.com/novelhive/proto/library/v1"
	notificationv1 "github.com/novelhive/proto/notification/v1"
	novelv1 "github.com/novelhive/proto/novel/v1"
	...
	type Clients struct {
	Novel        novelv1.NovelServiceClient
	User         userv1.UserServiceClient
	Comment      commentv1.CommentServiceClient
	Library      libraryv1.LibraryServiceClient
	Notification notificationv1.NotificationServiceClient

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
		conns:        []*grpc.ClientConn{userConn, novelConn, commentConn, libraryConn, notificationConn},
	}
	}

	}
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

