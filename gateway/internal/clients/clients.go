package clients

import (
	"log"

	commentv1 "github.com/novelhive/proto/comment/v1"
	libraryv1 "github.com/novelhive/proto/library/v1"
	novelv1 "github.com/novelhive/proto/novel/v1"
	userv1 "github.com/novelhive/proto/user/v1"
	"github.com/novelhive/pkg/grpcauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients holds all gRPC service clients.
type Clients struct {
	Novel   novelv1.NovelServiceClient
	User    userv1.UserServiceClient
	Comment commentv1.CommentServiceClient
	Library libraryv1.LibraryServiceClient

	conns []*grpc.ClientConn
}

// New creates gRPC clients to all downstream microservices.
// apiKey is injected as x-internal-key metadata on every call so that
// services can reject calls that don't originate from the gateway.
func New(userAddr, novelAddr, commentAddr, libraryAddr, apiKey string) *Clients {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(grpcauth.NewCredentials(apiKey)),
	}

	userConn := mustDial(userAddr, opts)
	novelConn := mustDial(novelAddr, opts)
	commentConn := mustDial(commentAddr, opts)
	libraryConn := mustDial(libraryAddr, opts)

	return &Clients{
		User:    userv1.NewUserServiceClient(userConn),
		Novel:   novelv1.NewNovelServiceClient(novelConn),
		Comment: commentv1.NewCommentServiceClient(commentConn),
		Library: libraryv1.NewLibraryServiceClient(libraryConn),
		conns:   []*grpc.ClientConn{userConn, novelConn, commentConn, libraryConn},
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

