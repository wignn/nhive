package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/novelhive/novel-service/internal/domain"
	"github.com/novelhive/novel-service/internal/events"
	"github.com/novelhive/novel-service/internal/repository"
)

type NovelUsecase struct {
	novelRepo   domain.NovelRepository
	chapterRepo domain.ChapterRepository
	genreRepo   domain.GenreRepository
	cache       *repository.RedisCache
	publisher   *events.NATSPublisher
}

func NewNovelUsecase(nr domain.NovelRepository, cr domain.ChapterRepository, gr domain.GenreRepository, c *repository.RedisCache, p *events.NATSPublisher) *NovelUsecase {
	return &NovelUsecase{novelRepo: nr, chapterRepo: cr, genreRepo: gr, cache: c, publisher: p}
}

func (uc *NovelUsecase) CreateNovel(title, synopsis, coverURL, author string, genreIDs []int) (*domain.Novel, error) {
	if strings.TrimSpace(title) == "" {
		return nil, domain.ErrInvalidInput
	}
	novel := &domain.Novel{
		ID: generateID(), Title: title, Slug: slug.Make(title),
		Synopsis: synopsis, CoverURL: coverURL, Author: author,
		Status: "ongoing", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := uc.novelRepo.Create(novel); err != nil {
		return nil, err
	}
	if len(genreIDs) > 0 {
		uc.novelRepo.SetGenres(novel.ID, genreIDs)
	}
	if uc.publisher != nil {
		uc.publisher.PublishNovelCreated(novel)
	}
	return novel, nil
}

func (uc *NovelUsecase) GetNovel(novelSlug string) (*domain.Novel, error) {
	if uc.cache != nil {
		if c, err := uc.cache.GetNovel(novelSlug); err == nil {
			return c, nil
		}
	}
	novel, err := uc.novelRepo.GetBySlug(novelSlug)
	if err != nil {
		return nil, err
	}
	if uc.cache != nil {
		uc.cache.SetNovel(novel)
	}
	return novel, nil
}

func (uc *NovelUsecase) ListNovels(p domain.ListNovelsParams) ([]*domain.Novel, int, error) {
	return uc.novelRepo.List(p)
}

func (uc *NovelUsecase) CreateChapter(novelID string, number int, title, content string) (*domain.Chapter, error) {
	if strings.TrimSpace(content) == "" {
		return nil, domain.ErrInvalidInput
	}
	ch := &domain.Chapter{
		ID: generateID(), NovelID: novelID, Number: number,
		Title: title, Content: content, WordCount: len(strings.Fields(content)),
		CreatedAt: time.Now(),
	}
	if err := uc.chapterRepo.Create(ch); err != nil {
		return nil, err
	}
	count, _ := uc.chapterRepo.CountByNovelID(novelID)
	if novel, err := uc.novelRepo.GetByID(novelID); err == nil {
		novel.TotalChapters = count
		novel.UpdatedAt = time.Now()
		uc.novelRepo.Update(novel)
		if uc.cache != nil {
			uc.cache.InvalidateNovel(novel.Slug)
		}
		if uc.publisher != nil {
			uc.publisher.PublishChapterPublished(novel, ch)
		}
	}
	return ch, nil
}

func (uc *NovelUsecase) GetChapter(novelSlug string, number int) (*domain.Chapter, error) {
	if uc.cache != nil {
		if c, err := uc.cache.GetChapter(novelSlug, number); err == nil {
			return c, nil
		}
	}
	ch, err := uc.chapterRepo.GetByNovelAndNumber(novelSlug, number)
	if err != nil {
		return nil, err
	}
	if uc.cache != nil {
		uc.cache.SetChapter(novelSlug, ch)
	}
	return ch, nil
}

func (uc *NovelUsecase) ListChapters(p domain.ListChaptersParams) ([]domain.ChapterSummary, int, error) {
	return uc.chapterRepo.List(p)
}

func (uc *NovelUsecase) ListGenres() ([]domain.Genre, error) {
	return uc.genreRepo.List()
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
