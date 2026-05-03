package domain

import (
	"time"
)

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Payload   map[string]interface{} `json:"payload"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
}

type NotificationRepository interface {
	Create(notification *Notification) error
	ListByUser(userID string, page, pageSize int) ([]*Notification, int, error)
	MarkAsRead(id, userID string) error
	SaveFCMToken(userID, token string) error
	GetFCMTokens(userID string) ([]string, error)
}
