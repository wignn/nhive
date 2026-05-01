package domain

import "time"

type Novel struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Synopsis      string    `json:"synopsis"`
	CoverURL      string    `json:"cover_url"`
	Author        string    `json:"author"`
	Status        string    `json:"status"`
	TotalChapters int       `json:"total_chapters"`
	Genres        []Genre   `json:"genres"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Chapter struct {
	ID        string    `json:"id"`
	NovelID   string    `json:"novel_id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	WordCount int       `json:"word_count"`
	CreatedAt time.Time `json:"created_at"`
}

type ChapterSummary struct {
	ID        string    `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	WordCount int       `json:"word_count"`
	CreatedAt time.Time `json:"created_at"`
}

type ListNovelsParams struct {
	Page      int
	PageSize  int
	GenreSlug string
	Status    string
	SortBy    string
}

type ListChaptersParams struct {
	NovelID  string
	Page     int
	PageSize int
}

type NovelRepository interface {
	Create(novel *Novel) error
	GetBySlug(slug string) (*Novel, error)
	GetByID(id string) (*Novel, error)
	List(params ListNovelsParams) ([]*Novel, int, error)
	Update(novel *Novel) error
	Delete(id string) error
	SetGenres(novelID string, genreIDs []int) error
}

type ChapterRepository interface {
	Create(chapter *Chapter) error
	GetByNovelAndNumber(novelSlug string, number int) (*Chapter, error)
	List(params ListChaptersParams) ([]ChapterSummary, int, error)
	Update(chapter *Chapter) error
	CountByNovelID(novelID string) (int, error)
}

type GenreRepository interface {
	List() ([]Genre, error)
	GetByID(id int) (*Genre, error)
	Create(name, slug string) (*Genre, error)
	Delete(id int) error
}
