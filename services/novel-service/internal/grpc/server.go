package grpcserver

import (
	"context"

	"github.com/novelhive/novel-service/internal/domain"
	"github.com/novelhive/novel-service/internal/usecase"
	novelv1 "github.com/novelhive/proto/novel/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NovelServiceServer struct {
	novelv1.UnimplementedNovelServiceServer
	uc *usecase.NovelUsecase
}

func NewNovelServiceServer(uc *usecase.NovelUsecase) *NovelServiceServer {
	return &NovelServiceServer{uc: uc}
}

// --- Novel RPCs ---

func (s *NovelServiceServer) GetNovel(ctx context.Context, req *novelv1.GetNovelRequest) (*novelv1.Novel, error) {
	novel, err := s.uc.GetNovel(req.Slug)
	if err != nil {
		if err == domain.ErrNovelNotFound {
			return nil, status.Error(codes.NotFound, "novel not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoNovel(novel), nil
}

func (s *NovelServiceServer) ListNovels(ctx context.Context, req *novelv1.ListNovelsRequest) (*novelv1.ListNovelsResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	novels, total, err := s.uc.ListNovels(domain.ListNovelsParams{
		Page:      page,
		PageSize:  pageSize,
		GenreSlug: req.GenreSlug,
		Status:    req.Status,
		SortBy:    req.SortBy,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var protoNovels []*novelv1.Novel
	for _, n := range novels {
		protoNovels = append(protoNovels, toProtoNovel(n))
	}
	return &novelv1.ListNovelsResponse{
		Novels:   protoNovels,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

func (s *NovelServiceServer) CreateNovel(ctx context.Context, req *novelv1.CreateNovelRequest) (*novelv1.Novel, error) {
	genreIDs := make([]int, len(req.GenreIds))
	for i, id := range req.GenreIds {
		genreIDs[i] = int(id)
	}
	novel, err := s.uc.CreateNovel(req.Title, req.Synopsis, req.CoverUrl, req.Author, genreIDs)
	if err != nil {
		if err == domain.ErrInvalidInput {
			return nil, status.Error(codes.InvalidArgument, "invalid input")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoNovel(novel), nil
}

func (s *NovelServiceServer) UpdateNovel(ctx context.Context, req *novelv1.UpdateNovelRequest) (*novelv1.Novel, error) {
	novel, err := s.uc.GetNovelByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "novel not found")
	}
	if req.Title != "" {
		novel.Title = req.Title
	}
	if req.Synopsis != "" {
		novel.Synopsis = req.Synopsis
	}
	if req.CoverUrl != "" {
		novel.CoverURL = req.CoverUrl
	}
	if req.Author != "" {
		novel.Author = req.Author
	}
	if req.Status != "" {
		novel.Status = req.Status
	}
	if err := s.uc.UpdateNovel(novel); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(req.GenreIds) > 0 {
		genreIDs := make([]int, len(req.GenreIds))
		for i, id := range req.GenreIds {
			genreIDs[i] = int(id)
		}
		s.uc.SetGenres(novel.ID, genreIDs)
	}
	updated, _ := s.uc.GetNovelByID(req.Id)
	return toProtoNovel(updated), nil
}

func (s *NovelServiceServer) DeleteNovel(ctx context.Context, req *novelv1.DeleteNovelRequest) (*novelv1.DeleteNovelResponse, error) {
	if err := s.uc.DeleteNovel(req.Id); err != nil {
		return nil, status.Error(codes.NotFound, "novel not found")
	}
	return &novelv1.DeleteNovelResponse{Success: true}, nil
}

// --- Chapter RPCs ---

func (s *NovelServiceServer) GetChapter(ctx context.Context, req *novelv1.GetChapterRequest) (*novelv1.Chapter, error) {
	ch, err := s.uc.GetChapter(req.NovelSlug, int(req.ChapterNumber))
	if err != nil {
		if err == domain.ErrChapterNotFound {
			return nil, status.Error(codes.NotFound, "chapter not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	novel, _ := s.uc.GetNovel(req.NovelSlug)
	protoChapter := toProtoChapter(ch)
	if novel != nil {
		_ = novel // total_chapters available if needed
	}
	return protoChapter, nil
}

func (s *NovelServiceServer) ListChapters(ctx context.Context, req *novelv1.ListChaptersRequest) (*novelv1.ListChaptersResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}
	summaries, total, err := s.uc.ListChapters(domain.ListChaptersParams{
		NovelID:  req.NovelId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var protoSummaries []*novelv1.ChapterSummary
	for _, cs := range summaries {
		protoSummaries = append(protoSummaries, &novelv1.ChapterSummary{
			Id:        cs.ID,
			Number:    int32(cs.Number),
			Title:     cs.Title,
			WordCount: int32(cs.WordCount),
			CreatedAt: cs.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return &novelv1.ListChaptersResponse{
		Chapters: protoSummaries,
		Total:    int32(total),
	}, nil
}

func (s *NovelServiceServer) CreateChapter(ctx context.Context, req *novelv1.CreateChapterRequest) (*novelv1.Chapter, error) {
	ch, err := s.uc.CreateChapter(req.NovelId, int(req.Number), req.Title, req.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toProtoChapter(ch), nil
}

func (s *NovelServiceServer) UpdateChapter(ctx context.Context, req *novelv1.UpdateChapterRequest) (*novelv1.Chapter, error) {
	ch, err := s.uc.UpdateChapter(req.Id, req.Title, req.Content)
	if err != nil {
		return nil, status.Error(codes.NotFound, "chapter not found")
	}
	return toProtoChapter(ch), nil
}

func (s *NovelServiceServer) DeleteChapter(ctx context.Context, req *novelv1.DeleteChapterRequest) (*novelv1.DeleteChapterResponse, error) {
	if err := s.uc.DeleteChapter(req.Id); err != nil {
		return nil, status.Error(codes.NotFound, "chapter not found")
	}
	return &novelv1.DeleteChapterResponse{Success: true}, nil
}

func (s *NovelServiceServer) ListGenres(ctx context.Context, req *novelv1.ListGenresRequest) (*novelv1.ListGenresResponse, error) {
	genres, err := s.uc.ListGenres()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var protoGenres []*novelv1.Genre
	for _, g := range genres {
		protoGenres = append(protoGenres, &novelv1.Genre{
			Id:   int32(g.ID),
			Name: g.Name,
			Slug: g.Slug,
		})
	}
	return &novelv1.ListGenresResponse{Genres: protoGenres}, nil
}

func (s *NovelServiceServer) CreateGenre(ctx context.Context, req *novelv1.CreateGenreRequest) (*novelv1.Genre, error) {
	g, err := s.uc.CreateGenre(req.Name)
	if err != nil {
		if err == domain.ErrInvalidInput {
			return nil, status.Error(codes.InvalidArgument, "invalid genre name")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &novelv1.Genre{
		Id:   int32(g.ID),
		Name: g.Name,
		Slug: g.Slug,
	}, nil
}

func (s *NovelServiceServer) DeleteGenre(ctx context.Context, req *novelv1.DeleteGenreRequest) (*novelv1.DeleteGenreResponse, error) {
	if err := s.uc.DeleteGenre(int(req.Id)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &novelv1.DeleteGenreResponse{Success: true}, nil
}


func toProtoNovel(n *domain.Novel) *novelv1.Novel {
	var genres []*novelv1.Genre
	for _, g := range n.Genres {
		genres = append(genres, &novelv1.Genre{
			Id:   int32(g.ID),
			Name: g.Name,
			Slug: g.Slug,
		})
	}
	return &novelv1.Novel{
		Id:            n.ID,
		Title:         n.Title,
		Slug:          n.Slug,
		Synopsis:      n.Synopsis,
		CoverUrl:      n.CoverURL,
		Author:        n.Author,
		Status:        n.Status,
		TotalChapters: int32(n.TotalChapters),
		Genres:        genres,
		CreatedAt:     n.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     n.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toProtoChapter(ch *domain.Chapter) *novelv1.Chapter {
	return &novelv1.Chapter{
		Id:        ch.ID,
		NovelId:   ch.NovelID,
		Number:    int32(ch.Number),
		Title:     ch.Title,
		Content:   ch.Content,
		WordCount: int32(ch.WordCount),
		CreatedAt: ch.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
