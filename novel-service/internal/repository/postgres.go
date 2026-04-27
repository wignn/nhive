package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/novel-service/internal/domain"
)

type PostgresNovelRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresNovelRepo(pool *pgxpool.Pool) *PostgresNovelRepo {
	return &PostgresNovelRepo{pool: pool}
}

func (r *PostgresNovelRepo) Create(novel *domain.Novel) error {
	query := `
		INSERT INTO novels (id, title, slug, synopsis, cover_url, author, status, total_chapters, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(context.Background(), query,
		novel.ID, novel.Title, novel.Slug, novel.Synopsis, novel.CoverURL,
		novel.Author, novel.Status, novel.TotalChapters, novel.CreatedAt, novel.UpdatedAt,
	)
	return err
}

func (r *PostgresNovelRepo) GetBySlug(slug string) (*domain.Novel, error) {
	query := `SELECT id, title, slug, synopsis, cover_url, author, status, total_chapters, created_at, updated_at FROM novels WHERE slug = $1`
	novel, err := r.scanNovel(r.pool.QueryRow(context.Background(), query, slug))
	if err != nil {
		return nil, err
	}
	// Load genres
	genres, err := r.getGenresForNovel(novel.ID)
	if err == nil {
		novel.Genres = genres
	}
	return novel, nil
}

func (r *PostgresNovelRepo) GetByID(id string) (*domain.Novel, error) {
	query := `SELECT id, title, slug, synopsis, cover_url, author, status, total_chapters, created_at, updated_at FROM novels WHERE id = $1`
	novel, err := r.scanNovel(r.pool.QueryRow(context.Background(), query, id))
	if err != nil {
		return nil, err
	}
	genres, err := r.getGenresForNovel(novel.ID)
	if err == nil {
		novel.Genres = genres
	}
	return novel, nil
}

func (r *PostgresNovelRepo) List(params domain.ListNovelsParams) ([]*domain.Novel, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 50 {
		params.PageSize = 20
	}

	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if params.Status != "" {
		where += fmt.Sprintf(" AND n.status = $%d", argIdx)
		args = append(args, params.Status)
		argIdx++
	}
	if params.GenreSlug != "" {
		where += fmt.Sprintf(` AND EXISTS (
			SELECT 1 FROM novel_genres ng JOIN genres g ON ng.genre_id = g.id
			WHERE ng.novel_id = n.id AND g.slug = $%d
		)`, argIdx)
		args = append(args, params.GenreSlug)
		argIdx++
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM novels n %s", where)
	var total int
	r.pool.QueryRow(context.Background(), countQuery, args...).Scan(&total)

	// Sort
	orderBy := "ORDER BY n.updated_at DESC"
	if params.SortBy == "title" {
		orderBy = "ORDER BY n.title ASC"
	} else if params.SortBy == "created" {
		orderBy = "ORDER BY n.created_at DESC"
	}

	offset := (params.Page - 1) * params.PageSize
	query := fmt.Sprintf(`
		SELECT n.id, n.title, n.slug, n.synopsis, n.cover_url, n.author, n.status, n.total_chapters, n.created_at, n.updated_at
		FROM novels n %s %s LIMIT $%d OFFSET $%d
	`, where, orderBy, argIdx, argIdx+1)
	args = append(args, params.PageSize, offset)

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var novels []*domain.Novel
	for rows.Next() {
		n := &domain.Novel{}
		if err := rows.Scan(&n.ID, &n.Title, &n.Slug, &n.Synopsis, &n.CoverURL,
			&n.Author, &n.Status, &n.TotalChapters, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, 0, err
		}
		genres, _ := r.getGenresForNovel(n.ID)
		n.Genres = genres
		novels = append(novels, n)
	}
	return novels, total, nil
}

func (r *PostgresNovelRepo) Update(novel *domain.Novel) error {
	query := `UPDATE novels SET title=$1, synopsis=$2, cover_url=$3, author=$4, status=$5, total_chapters=$6, updated_at=$7 WHERE id=$8`
	_, err := r.pool.Exec(context.Background(), query,
		novel.Title, novel.Synopsis, novel.CoverURL, novel.Author, novel.Status, novel.TotalChapters, novel.UpdatedAt, novel.ID,
	)
	return err
}

func (r *PostgresNovelRepo) Delete(id string) error {
	_, err := r.pool.Exec(context.Background(), "DELETE FROM novels WHERE id = $1", id)
	return err
}

func (r *PostgresNovelRepo) SetGenres(novelID string, genreIDs []int) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	tx.Exec(context.Background(), "DELETE FROM novel_genres WHERE novel_id = $1", novelID)
	for _, gid := range genreIDs {
		tx.Exec(context.Background(), "INSERT INTO novel_genres (novel_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", novelID, gid)
	}
	return tx.Commit(context.Background())
}

func (r *PostgresNovelRepo) getGenresForNovel(novelID string) ([]domain.Genre, error) {
	query := `SELECT g.id, g.name, g.slug FROM genres g JOIN novel_genres ng ON g.id = ng.genre_id WHERE ng.novel_id = $1`
	rows, err := r.pool.Query(context.Background(), query, novelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var genres []domain.Genre
	for rows.Next() {
		var g domain.Genre
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *PostgresNovelRepo) scanNovel(row pgx.Row) (*domain.Novel, error) {
	var n domain.Novel
	err := row.Scan(&n.ID, &n.Title, &n.Slug, &n.Synopsis, &n.CoverURL,
		&n.Author, &n.Status, &n.TotalChapters, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNovelNotFound
		}
		return nil, err
	}
	return &n, nil
}

// Chapter repository
type PostgresChapterRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresChapterRepo(pool *pgxpool.Pool) *PostgresChapterRepo {
	return &PostgresChapterRepo{pool: pool}
}

func (r *PostgresChapterRepo) Create(ch *domain.Chapter) error {
	query := `INSERT INTO chapters (id, novel_id, number, title, content, word_count, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(context.Background(), query,
		ch.ID, ch.NovelID, ch.Number, ch.Title, ch.Content, ch.WordCount, ch.CreatedAt,
	)
	return err
}

func (r *PostgresChapterRepo) GetByNovelAndNumber(novelSlug string, number int) (*domain.Chapter, error) {
	query := `
		SELECT c.id, c.novel_id, c.number, c.title, c.content, c.word_count, c.created_at
		FROM chapters c JOIN novels n ON c.novel_id = n.id
		WHERE n.slug = $1 AND c.number = $2
	`
	var ch domain.Chapter
	err := r.pool.QueryRow(context.Background(), query, novelSlug, number).Scan(
		&ch.ID, &ch.NovelID, &ch.Number, &ch.Title, &ch.Content, &ch.WordCount, &ch.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChapterNotFound
		}
		return nil, err
	}
	return &ch, nil
}

func (r *PostgresChapterRepo) List(params domain.ListChaptersParams) ([]domain.ChapterSummary, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 50
	}
	var total int
	r.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM chapters WHERE novel_id = $1", params.NovelID).Scan(&total)

	offset := (params.Page - 1) * params.PageSize
	query := `SELECT id, number, title, word_count, created_at FROM chapters WHERE novel_id = $1 ORDER BY number ASC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(context.Background(), query, params.NovelID, params.PageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var chapters []domain.ChapterSummary
	for rows.Next() {
		var cs domain.ChapterSummary
		if err := rows.Scan(&cs.ID, &cs.Number, &cs.Title, &cs.WordCount, &cs.CreatedAt); err != nil {
			return nil, 0, err
		}
		chapters = append(chapters, cs)
	}
	return chapters, total, nil
}

func (r *PostgresChapterRepo) Update(ch *domain.Chapter) error {
	query := `UPDATE chapters SET title=$1, content=$2, word_count=$3 WHERE id=$4`
	_, err := r.pool.Exec(context.Background(), query, ch.Title, ch.Content, ch.WordCount, ch.ID)
	return err
}

func (r *PostgresChapterRepo) CountByNovelID(novelID string) (int, error) {
	var count int
	err := r.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM chapters WHERE novel_id = $1", novelID).Scan(&count)
	return count, err
}

// Genre repository
type PostgresGenreRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresGenreRepo(pool *pgxpool.Pool) *PostgresGenreRepo {
	return &PostgresGenreRepo{pool: pool}
}

func (r *PostgresGenreRepo) List() ([]domain.Genre, error) {
	rows, err := r.pool.Query(context.Background(), "SELECT id, name, slug FROM genres ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var genres []domain.Genre
	for rows.Next() {
		var g domain.Genre
		rows.Scan(&g.ID, &g.Name, &g.Slug)
		genres = append(genres, g)
	}
	return genres, nil
}

func (r *PostgresGenreRepo) GetByID(id int) (*domain.Genre, error) {
	var g domain.Genre
	err := r.pool.QueryRow(context.Background(), "SELECT id, name, slug FROM genres WHERE id = $1", id).Scan(&g.ID, &g.Name, &g.Slug)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
