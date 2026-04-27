package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/library-service/internal/domain"
)

type PostgresLibraryRepo struct{ pool *pgxpool.Pool }
func NewPostgresLibraryRepo(pool *pgxpool.Pool) *PostgresLibraryRepo { return &PostgresLibraryRepo{pool: pool} }

func (r *PostgresLibraryRepo) AddToLibrary(e *domain.LibraryEntry) error {
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO reading_lists (id, user_id, novel_id, status, created_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (user_id, novel_id) DO NOTHING`,
		e.ID, e.UserID, e.NovelID, e.Status, e.CreatedAt)
	return err
}

func (r *PostgresLibraryRepo) GetLibrary(userID, status string, page, pageSize int) ([]*domain.LibraryEntry, int, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 50 { pageSize = 20 }
	where := "WHERE user_id = $1"
	args := []interface{}{userID}
	if status != "" {
		where += " AND status = $2"
		args = append(args, status)
	}
	var total int
	r.pool.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM reading_lists %s", where), args...).Scan(&total)

	offset := (page - 1) * pageSize
	q := fmt.Sprintf("SELECT id, user_id, novel_id, status, created_at FROM reading_lists %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)
	rows, err := r.pool.Query(context.Background(), q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var entries []*domain.LibraryEntry
	for rows.Next() {
		e := &domain.LibraryEntry{}
		rows.Scan(&e.ID, &e.UserID, &e.NovelID, &e.Status, &e.CreatedAt)
		entries = append(entries, e)
	}
	return entries, total, nil
}

func (r *PostgresLibraryRepo) UpdateStatus(userID, novelID, status string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE reading_lists SET status=$1 WHERE user_id=$2 AND novel_id=$3", status, userID, novelID)
	return err
}

func (r *PostgresLibraryRepo) RemoveFromLibrary(userID, novelID string) error {
	_, err := r.pool.Exec(context.Background(), "DELETE FROM reading_lists WHERE user_id=$1 AND novel_id=$2", userID, novelID)
	return err
}

// Bookmark repo
type PostgresBookmarkRepo struct{ pool *pgxpool.Pool }
func NewPostgresBookmarkRepo(pool *pgxpool.Pool) *PostgresBookmarkRepo { return &PostgresBookmarkRepo{pool: pool} }

func (r *PostgresBookmarkRepo) Add(b *domain.Bookmark) error {
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO bookmarks (id, user_id, novel_id, chapter_id, note, created_at) VALUES ($1,$2,$3,$4,$5,$6)",
		b.ID, b.UserID, b.NovelID, b.ChapterID, b.Note, b.CreatedAt)
	return err
}

func (r *PostgresBookmarkRepo) List(userID, novelID string) ([]*domain.Bookmark, error) {
	q := "SELECT id, user_id, novel_id, chapter_id, note, created_at FROM bookmarks WHERE user_id=$1"
	args := []interface{}{userID}
	if novelID != "" {
		q += " AND novel_id=$2"
		args = append(args, novelID)
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.pool.Query(context.Background(), q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var bookmarks []*domain.Bookmark
	for rows.Next() {
		b := &domain.Bookmark{}
		rows.Scan(&b.ID, &b.UserID, &b.NovelID, &b.ChapterID, &b.Note, &b.CreatedAt)
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func (r *PostgresBookmarkRepo) Remove(bookmarkID, userID string) error {
	_, err := r.pool.Exec(context.Background(), "DELETE FROM bookmarks WHERE id=$1 AND user_id=$2", bookmarkID, userID)
	return err
}

// Progress repo
type PostgresProgressRepo struct{ pool *pgxpool.Pool }
func NewPostgresProgressRepo(pool *pgxpool.Pool) *PostgresProgressRepo { return &PostgresProgressRepo{pool: pool} }

func (r *PostgresProgressRepo) Get(userID, novelID string) (*domain.ReadingProgress, error) {
	var p domain.ReadingProgress
	err := r.pool.QueryRow(context.Background(),
		"SELECT user_id, novel_id, chapter_number, scroll_position, updated_at FROM reading_progress WHERE user_id=$1 AND novel_id=$2",
		userID, novelID).Scan(&p.UserID, &p.NovelID, &p.ChapterNumber, &p.ScrollPosition, &p.UpdatedAt)
	if err != nil { return nil, err }
	return &p, nil
}

func (r *PostgresProgressRepo) Save(p *domain.ReadingProgress) error {
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO reading_progress (user_id, novel_id, chapter_number, scroll_position, updated_at)
		VALUES ($1,$2,$3,$4,$5) ON CONFLICT (user_id, novel_id) DO UPDATE SET chapter_number=$3, scroll_position=$4, updated_at=$5`,
		p.UserID, p.NovelID, p.ChapterNumber, p.ScrollPosition, p.UpdatedAt)
	return err
}
