package domain

import "time"

type Comment struct {
	ID         string     `json:"id"`
	ChapterID  string     `json:"chapter_id"`
	UserID     string     `json:"user_id"`
	Username   string     `json:"username"`
	AvatarURL  string     `json:"avatar_url"`
	Content    string     `json:"content"`
	ParentID   string     `json:"parent_id,omitempty"`
	Path       string     `json:"path"` // ltree path
	LikesCount int        `json:"likes_count"`
	ReplyCount int        `json:"reply_count"`
	Replies    []*Comment `json:"replies,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id string) (*Comment, error)
	ListByChapter(chapterID string, page, pageSize int, sortBy string) ([]*Comment, int, error)
	Delete(id string) error
	IncrementLikes(id string) error
	DecrementLikes(id string) error
}

type CommentLikeRepository interface {
	HasLiked(commentID, userID string) (bool, error)
	Like(commentID, userID string) error
	Unlike(commentID, userID string) error
}
