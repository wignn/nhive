package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/user-service/internal/domain"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, avatar_url, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(context.Background(), query,
		user.ID, user.Username, user.Email, user.PasswordHash,
		user.AvatarURL, user.Role, user.CreatedAt,
	)
	return err
}

func (r *PostgresUserRepo) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, created_at FROM users WHERE id = $1`
	return r.scanUser(r.pool.QueryRow(context.Background(), query, id))
}

func (r *PostgresUserRepo) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, created_at FROM users WHERE email = $1`
	return r.scanUser(r.pool.QueryRow(context.Background(), query, email))
}

func (r *PostgresUserRepo) GetByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, avatar_url, role, created_at FROM users WHERE username = $1`
	return r.scanUser(r.pool.QueryRow(context.Background(), query, username))
}

func (r *PostgresUserRepo) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.pool.QueryRow(context.Background(), query, email).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepo) ExistsByUsername(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.pool.QueryRow(context.Background(), query, username).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepo) scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
