package domain

import "time"

type LibraryEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	NovelID   string    `json:"novel_id"`
	Status    string    `json:"status"` // reading, completed, plan_to_read, dropped
	CreatedAt time.Time `json:"created_at"`
}

type Bookmark struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	NovelID   string    `json:"novel_id"`
	ChapterID string    `json:"chapter_id"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadingProgress struct {
	UserID         string    `json:"user_id"`
	NovelID        string    `json:"novel_id"`
	ChapterNumber  int       `json:"chapter_number"`
	ScrollPosition float64   `json:"scroll_position"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LibraryRepository interface {
	AddToLibrary(entry *LibraryEntry) error
	GetLibrary(userID, status string, page, pageSize int) ([]*LibraryEntry, int, error)
	UpdateStatus(userID, novelID, status string) error
	RemoveFromLibrary(userID, novelID string) error
	GetUsersByNovel(novelID string) ([]string, error)
}

type BookmarkRepository interface {
	Add(bookmark *Bookmark) error
	List(userID, novelID string) ([]*Bookmark, error)
	Remove(bookmarkID, userID string) error
}

type ProgressRepository interface {
	Get(userID, novelID string) (*ReadingProgress, error)
	Save(progress *ReadingProgress) error
}
