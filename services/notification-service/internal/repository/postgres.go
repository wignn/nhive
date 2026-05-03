package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelhive/notification-service/internal/domain"
)

type PostgresNotificationRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificationRepo(pool *pgxpool.Pool) *PostgresNotificationRepo {
	return &PostgresNotificationRepo{pool: pool}
}

func (r *PostgresNotificationRepo) Create(n *domain.Notification) error {
	payloadJSON, _ := json.Marshal(n.Payload)
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO notifications (id, user_id, type, title, body, payload, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		n.ID, n.UserID, n.Type, n.Title, n.Body, payloadJSON, n.IsRead, n.CreatedAt)
	return err
}

func (r *PostgresNotificationRepo) ListByUser(userID string, page, pageSize int) ([]*domain.Notification, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var total int
	r.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM notifications WHERE user_id = $1", userID).Scan(&total)

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(context.Background(),
		`SELECT id, user_id, type, title, body, payload, is_read, created_at
		FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		var payloadJSON []byte
		err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &payloadJSON, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		json.Unmarshal(payloadJSON, &n.Payload)
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (r *PostgresNotificationRepo) MarkAsRead(id, userID string) error {
	_, err := r.pool.Exec(context.Background(),
		"UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2",
		id, userID)
	return err
}

func (r *PostgresNotificationRepo) SaveFCMToken(userID, token string) error {
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO fcm_tokens (user_id, token, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, token) DO UPDATE SET updated_at = CURRENT_TIMESTAMP`,
		userID, token)
	return err
}

func (r *PostgresNotificationRepo) GetFCMTokens(userID string) ([]string, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT token FROM fcm_tokens WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
