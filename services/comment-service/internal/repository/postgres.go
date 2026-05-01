package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/comment-service/internal/domain"
)

type PostgresCommentRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresCommentRepo(pool *pgxpool.Pool) *PostgresCommentRepo {
	return &PostgresCommentRepo{pool: pool}
}

func (r *PostgresCommentRepo) Create(c *domain.Comment) error {
	query := `INSERT INTO comments (id, chapter_id, user_id, username, avatar_url, content, parent_id, path, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::ltree, $9)`
	_, err := r.pool.Exec(context.Background(), query,
		c.ID, c.ChapterID, c.UserID, c.Username, c.AvatarURL, c.Content, nullStr(c.ParentID), c.Path, c.CreatedAt,
	)
	return err
}

func (r *PostgresCommentRepo) GetByID(id string) (*domain.Comment, error) {
	query := `SELECT id, chapter_id, user_id, username, avatar_url, content, COALESCE(parent_id::text, ''), path::text, likes_count, created_at
		FROM comments WHERE id = $1 AND deleted_at IS NULL`
	var c domain.Comment
	err := r.pool.QueryRow(context.Background(), query, id).Scan(
		&c.ID, &c.ChapterID, &c.UserID, &c.Username, &c.AvatarURL, &c.Content, &c.ParentID, &c.Path, &c.LikesCount, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCommentNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PostgresCommentRepo) ListByChapter(chapterID string, page, pageSize int, sortBy string) ([]*domain.Comment, int, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 50 { pageSize = 20 }

	// Count total root comments
	var total int
	r.pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM comments WHERE chapter_id = $1 AND parent_id IS NULL AND deleted_at IS NULL",
		chapterID).Scan(&total)

	order := "created_at DESC"
	if sortBy == "likes" { order = "likes_count DESC" }
	if sortBy == "oldest" { order = "created_at ASC" }

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT id, chapter_id, user_id, username, avatar_url, content, '', path::text, likes_count, created_at
		FROM comments WHERE chapter_id = $1 AND parent_id IS NULL AND deleted_at IS NULL
		ORDER BY %s LIMIT $2 OFFSET $3`, order)

	rows, err := r.pool.Query(context.Background(), query, chapterID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{}
		rows.Scan(&c.ID, &c.ChapterID, &c.UserID, &c.Username, &c.AvatarURL, &c.Content, &c.ParentID, &c.Path, &c.LikesCount, &c.CreatedAt)
		// Load replies (1 level deep)
		c.Replies, c.ReplyCount = r.loadReplies(c.ID)
		comments = append(comments, c)
	}
	return comments, total, nil
}

func (r *PostgresCommentRepo) loadReplies(parentID string) ([]*domain.Comment, int) {
	query := `SELECT id, chapter_id, user_id, username, avatar_url, content, COALESCE(parent_id::text, ''), path::text, likes_count, created_at
		FROM comments WHERE parent_id = $1 AND deleted_at IS NULL ORDER BY created_at ASC LIMIT 5`
	rows, err := r.pool.Query(context.Background(), query, parentID)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()
	var replies []*domain.Comment
	for rows.Next() {
		c := &domain.Comment{}
		rows.Scan(&c.ID, &c.ChapterID, &c.UserID, &c.Username, &c.AvatarURL, &c.Content, &c.ParentID, &c.Path, &c.LikesCount, &c.CreatedAt)
		replies = append(replies, c)
	}
	var count int
	r.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM comments WHERE parent_id = $1 AND deleted_at IS NULL", parentID).Scan(&count)
	return replies, count
}

func (r *PostgresCommentRepo) Delete(id string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE comments SET deleted_at = NOW(), content = '[deleted]' WHERE id = $1", id)
	return err
}

func (r *PostgresCommentRepo) IncrementLikes(id string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE comments SET likes_count = likes_count + 1 WHERE id = $1", id)
	return err
}

func (r *PostgresCommentRepo) DecrementLikes(id string) error {
	_, err := r.pool.Exec(context.Background(), "UPDATE comments SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = $1", id)
	return err
}

func nullStr(s string) interface{} {
	if s == "" { return nil }
	return s
}

// Like repository
type PostgresLikeRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresLikeRepo(pool *pgxpool.Pool) *PostgresLikeRepo {
	return &PostgresLikeRepo{pool: pool}
}

func (r *PostgresLikeRepo) HasLiked(commentID, userID string) (bool, error) {
	var exists bool
	r.pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id=$1 AND user_id=$2)", commentID, userID).Scan(&exists)
	return exists, nil
}

func (r *PostgresLikeRepo) Like(commentID, userID string) error {
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO comment_likes (comment_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", commentID, userID)
	return err
}

func (r *PostgresLikeRepo) Unlike(commentID, userID string) error {
	_, err := r.pool.Exec(context.Background(),
		"DELETE FROM comment_likes WHERE comment_id=$1 AND user_id=$2", commentID, userID)
	return err
}
