package grpcserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/novelhive/comment-service/internal/domain"
	"github.com/novelhive/comment-service/internal/repository"
	commentv1 "github.com/novelhive/proto/comment/v1"
	userv1 "github.com/novelhive/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentServiceServer struct {
	commentv1.UnimplementedCommentServiceServer
	commentRepo *repository.PostgresCommentRepo
	likeRepo    *repository.PostgresLikeRepo
	userClient  userv1.UserServiceClient
	logger      *zap.Logger
}

func NewCommentServiceServer(
	commentRepo *repository.PostgresCommentRepo,
	likeRepo *repository.PostgresLikeRepo,
	userClient userv1.UserServiceClient,
	logger *zap.Logger,
) *CommentServiceServer {
	return &CommentServiceServer{
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		userClient:  userClient,
		logger:      logger,
	}
}

func (s *CommentServiceServer) ListComments(ctx context.Context, req *commentv1.ListCommentsRequest) (*commentv1.ListCommentsResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	comments, total, err := s.commentRepo.ListByChapter(req.ChapterId, page, pageSize, "newest")
	if err != nil {
		s.logger.Error("failed to list comments",
			zap.String("chapter_id", req.ChapterId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to list comments")
	}

	var pbComments []*commentv1.Comment
	for _, c := range comments {
		pbComments = append(pbComments, domainToProto(c))
	}

	return &commentv1.ListCommentsResponse{
		Comments: pbComments,
		Total:    int32(total),
		Page:     int32(page),
	}, nil
}

func (s *CommentServiceServer) CreateComment(ctx context.Context, req *commentv1.CreateCommentRequest) (*commentv1.Comment, error) {
	id := genID()
	now := time.Now()

	username, avatarURL := s.resolveUser(ctx, req.UserId)

	c := &domain.Comment{
		ID:        id,
		ChapterID: req.ChapterId,
		UserID:    req.UserId,
		Username:  username,
		AvatarURL: avatarURL,
		Content:   req.Content,
		Path:      id,
		CreatedAt: now,
	}

	if err := s.commentRepo.Create(c); err != nil {
		s.logger.Error("failed to create comment",
			zap.String("user_id", req.UserId),
			zap.String("chapter_id", req.ChapterId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to create comment")
	}

	s.logger.Info("comment created",
		zap.String("comment_id", id),
		zap.String("user_id", req.UserId),
		zap.String("chapter_id", req.ChapterId),
	)

	return domainToProto(c), nil
}

func (s *CommentServiceServer) ReplyToComment(ctx context.Context, req *commentv1.ReplyToCommentRequest) (*commentv1.Comment, error) {
	parent, err := s.commentRepo.GetByID(req.ParentId)
	if err != nil {
		s.logger.Warn("parent comment not found for reply",
			zap.String("parent_id", req.ParentId),
			zap.Error(err),
		)
		return nil, status.Error(codes.NotFound, "parent comment not found")
	}

	id := genID()
	now := time.Now()
	username, avatarURL := s.resolveUser(ctx, req.UserId)

	c := &domain.Comment{
		ID:        id,
		ChapterID: parent.ChapterID,
		UserID:    req.UserId,
		Username:  username,
		AvatarURL: avatarURL,
		Content:   req.Content,
		ParentID:  req.ParentId,
		Path:      parent.Path + "." + id,
		CreatedAt: now,
	}

	if err := s.commentRepo.Create(c); err != nil {
		s.logger.Error("failed to create reply",
			zap.String("parent_id", req.ParentId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to create reply")
	}

	s.logger.Info("reply created",
		zap.String("comment_id", id),
		zap.String("parent_id", req.ParentId),
		zap.String("user_id", req.UserId),
	)

	return domainToProto(c), nil
}

func (s *CommentServiceServer) LikeComment(ctx context.Context, req *commentv1.LikeCommentRequest) (*commentv1.LikeCommentResponse, error) {
	if err := s.likeRepo.Like(req.CommentId, req.UserId); err != nil {
		s.logger.Error("failed to like comment",
			zap.String("comment_id", req.CommentId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
	}

	if err := s.commentRepo.IncrementLikes(req.CommentId); err != nil {
		s.logger.Error("failed to increment likes count",
			zap.String("comment_id", req.CommentId),
			zap.Error(err),
		)
	}

	comment, _ := s.commentRepo.GetByID(req.CommentId)
	likesCount := int32(0)
	if comment != nil {
		likesCount = int32(comment.LikesCount)
	}

	return &commentv1.LikeCommentResponse{LikesCount: likesCount, Liked: true}, nil
}

func (s *CommentServiceServer) DeleteComment(ctx context.Context, req *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	if err := s.commentRepo.Delete(req.CommentId); err != nil {
		s.logger.Error("failed to delete comment",
			zap.String("comment_id", req.CommentId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return &commentv1.DeleteCommentResponse{Success: false}, nil
	}

	s.logger.Info("comment deleted",
		zap.String("comment_id", req.CommentId),
		zap.String("user_id", req.UserId),
	)

	return &commentv1.DeleteCommentResponse{Success: true}, nil
}

// resolveUser fetches username and avatar from user-service via gRPC.
func (s *CommentServiceServer) resolveUser(ctx context.Context, userID string) (string, string) {
	if s.userClient == nil {
		return "", ""
	}
	resp, err := s.userClient.GetProfile(ctx, &userv1.GetProfileRequest{UserId: userID})
	if err != nil {
		s.logger.Warn("failed to resolve user profile",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return "", ""
	}
	return resp.Username, resp.AvatarUrl
}

func domainToProto(c *domain.Comment) *commentv1.Comment {
	return &commentv1.Comment{
		Id:         c.ID,
		ChapterId:  c.ChapterID,
		UserId:     c.UserID,
		Username:   c.Username,
		AvatarUrl:  c.AvatarURL,
		Content:    c.Content,
		ParentId:   c.ParentID,
		LikesCount: int32(c.LikesCount),
		CreatedAt:  c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func genID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
