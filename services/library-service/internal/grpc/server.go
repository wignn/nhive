package grpcserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/novelhive/library-service/internal/domain"
	"github.com/novelhive/library-service/internal/repository"
	libraryv1 "github.com/novelhive/proto/library/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LibraryServiceServer struct {
	libraryv1.UnimplementedLibraryServiceServer
	libraryRepo  *repository.PostgresLibraryRepo
	bookmarkRepo *repository.PostgresBookmarkRepo
	progressRepo *repository.PostgresProgressRepo
	logger       *zap.Logger
}

func NewLibraryServiceServer(
	libraryRepo *repository.PostgresLibraryRepo,
	bookmarkRepo *repository.PostgresBookmarkRepo,
	progressRepo *repository.PostgresProgressRepo,
	logger *zap.Logger,
) *LibraryServiceServer {
	return &LibraryServiceServer{
		libraryRepo:  libraryRepo,
		bookmarkRepo: bookmarkRepo,
		progressRepo: progressRepo,
		logger:       logger,
	}
}

func (s *LibraryServiceServer) GetLibrary(ctx context.Context, req *libraryv1.GetLibraryRequest) (*libraryv1.GetLibraryResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	entries, total, err := s.libraryRepo.GetLibrary(req.UserId, req.Status, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get library",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to get library")
	}

	var pbEntries []*libraryv1.LibraryEntry
	for _, e := range entries {
		pbEntries = append(pbEntries, &libraryv1.LibraryEntry{
			Id:        e.ID,
			UserId:    e.UserID,
			NovelId:   e.NovelID,
			Status:    e.Status,
			CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &libraryv1.GetLibraryResponse{Entries: pbEntries, Total: int32(total)}, nil
}

func (s *LibraryServiceServer) AddToLibrary(ctx context.Context, req *libraryv1.AddToLibraryRequest) (*libraryv1.LibraryEntry, error) {
	id := genID()
	now := time.Now()
	entry := &domain.LibraryEntry{
		ID:        id,
		UserID:    req.UserId,
		NovelID:   req.NovelId,
		Status:    "reading",
		CreatedAt: now,
	}

	if err := s.libraryRepo.AddToLibrary(entry); err != nil {
		s.logger.Error("failed to add to library",
			zap.String("user_id", req.UserId),
			zap.String("novel_id", req.NovelId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to add to library")
	}

	s.logger.Info("novel added to library",
		zap.String("user_id", req.UserId),
		zap.String("novel_id", req.NovelId),
	)

	return &libraryv1.LibraryEntry{
		Id:        id,
		UserId:    req.UserId,
		NovelId:   req.NovelId,
		Status:    "reading",
		CreatedAt: now.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *LibraryServiceServer) UpdateStatus(ctx context.Context, req *libraryv1.UpdateStatusRequest) (*libraryv1.LibraryEntry, error) {
	if err := s.libraryRepo.UpdateStatus(req.UserId, req.NovelId, req.Status); err != nil {
		s.logger.Error("failed to update library status",
			zap.String("user_id", req.UserId),
			zap.String("novel_id", req.NovelId),
			zap.String("status", req.Status),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to update status")
	}

	return &libraryv1.LibraryEntry{
		UserId:  req.UserId,
		NovelId: req.NovelId,
		Status:  req.Status,
	}, nil
}

func (s *LibraryServiceServer) RemoveFromLibrary(ctx context.Context, req *libraryv1.RemoveFromLibraryRequest) (*libraryv1.RemoveResponse, error) {
	if err := s.libraryRepo.RemoveFromLibrary(req.UserId, req.NovelId); err != nil {
		s.logger.Error("failed to remove from library",
			zap.String("user_id", req.UserId),
			zap.String("novel_id", req.NovelId),
			zap.Error(err),
		)
		return &libraryv1.RemoveResponse{Success: false}, nil
	}

	s.logger.Info("novel removed from library",
		zap.String("user_id", req.UserId),
		zap.String("novel_id", req.NovelId),
	)

	return &libraryv1.RemoveResponse{Success: true}, nil
}

func (s *LibraryServiceServer) GetBookmarks(ctx context.Context, req *libraryv1.GetBookmarksRequest) (*libraryv1.GetBookmarksResponse, error) {
	bookmarks, err := s.bookmarkRepo.List(req.UserId, req.NovelId)
	if err != nil {
		s.logger.Error("failed to get bookmarks",
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to get bookmarks")
	}

	var pbBookmarks []*libraryv1.Bookmark
	for _, b := range bookmarks {
		pbBookmarks = append(pbBookmarks, &libraryv1.Bookmark{
			Id:        b.ID,
			UserId:    b.UserID,
			NovelId:   b.NovelID,
			ChapterId: b.ChapterID,
			Note:      b.Note,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &libraryv1.GetBookmarksResponse{Bookmarks: pbBookmarks}, nil
}

func (s *LibraryServiceServer) AddBookmark(ctx context.Context, req *libraryv1.AddBookmarkRequest) (*libraryv1.Bookmark, error) {
	id := genID()
	now := time.Now()
	b := &domain.Bookmark{
		ID:        id,
		UserID:    req.UserId,
		NovelID:   req.NovelId,
		ChapterID: req.ChapterId,
		Note:      req.Note,
		CreatedAt: now,
	}

	if err := s.bookmarkRepo.Add(b); err != nil {
		s.logger.Error("failed to add bookmark",
			zap.String("user_id", req.UserId),
			zap.String("novel_id", req.NovelId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to add bookmark")
	}

	return &libraryv1.Bookmark{
		Id:        id,
		UserId:    req.UserId,
		NovelId:   req.NovelId,
		ChapterId: req.ChapterId,
		Note:      req.Note,
		CreatedAt: now.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *LibraryServiceServer) RemoveBookmark(ctx context.Context, req *libraryv1.RemoveBookmarkRequest) (*libraryv1.RemoveResponse, error) {
	if err := s.bookmarkRepo.Remove(req.BookmarkId, req.UserId); err != nil {
		s.logger.Error("failed to remove bookmark",
			zap.String("bookmark_id", req.BookmarkId),
			zap.String("user_id", req.UserId),
			zap.Error(err),
		)
		return &libraryv1.RemoveResponse{Success: false}, nil
	}

	return &libraryv1.RemoveResponse{Success: true}, nil
}

func (s *LibraryServiceServer) GetProgress(ctx context.Context, req *libraryv1.GetProgressRequest) (*libraryv1.ReadingProgress, error) {
	p, err := s.progressRepo.Get(req.UserId, req.NovelId)
	if err != nil {
		// Not found is not an error — return default progress
		return &libraryv1.ReadingProgress{
			UserId:        req.UserId,
			NovelId:       req.NovelId,
			ChapterNumber: 1,
		}, nil
	}

	return &libraryv1.ReadingProgress{
		UserId:         p.UserID,
		NovelId:        p.NovelID,
		ChapterNumber:  int32(p.ChapterNumber),
		ScrollPosition: float32(p.ScrollPosition),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *LibraryServiceServer) SaveProgress(ctx context.Context, req *libraryv1.SaveProgressRequest) (*libraryv1.ReadingProgress, error) {
	now := time.Now()
	p := &domain.ReadingProgress{
		UserID:         req.UserId,
		NovelID:        req.NovelId,
		ChapterNumber:  int(req.ChapterNumber),
		ScrollPosition: float64(req.ScrollPosition),
		UpdatedAt:      now,
	}

	if err := s.progressRepo.Save(p); err != nil {
		s.logger.Error("failed to save reading progress",
			zap.String("user_id", req.UserId),
			zap.String("novel_id", req.NovelId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to save progress")
	}

	return &libraryv1.ReadingProgress{
		UserId:         req.UserId,
		NovelId:        req.NovelId,
		ChapterNumber:  req.ChapterNumber,
		ScrollPosition: req.ScrollPosition,
		UpdatedAt:      now.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func genID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
