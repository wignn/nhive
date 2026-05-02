package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)


type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarURL    string    `json:"avatar_url"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Novel struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Synopsis      string    `json:"synopsis"`
	CoverURL      string    `json:"cover_url"`
	Author        string    `json:"author"`
	Status        string    `json:"status"`
	TotalChapters int       `json:"total_chapters"`
	Genres        []string  `json:"genres"`
	ViewCount     int       `json:"view_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

type Comment struct {
	ID        string    `json:"id"`
	ChapterID string    `json:"chapter_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	ParentID  string    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type LibraryEntry struct {
	UserID     string `json:"user_id"`
	NovelID    string `json:"novel_id"`
	NovelTitle string `json:"novel_title"`
	NovelSlug  string `json:"novel_slug"`
	Status     string `json:"status"`
	Progress   int    `json:"progress"`
	Total      int    `json:"total"`
}

type ReadingProgress struct {
	UserID         string `json:"user_id"`
	NovelID        string `json:"novel_id"`
	ChapterNumber  int    `json:"chapter_number"`
	ScrollPosition int    `json:"scroll_position"`
}

type Store struct {
	mu       sync.RWMutex
	users    map[string]*User
	novels   map[string]*Novel
	chapters map[string][]*Chapter // keyed by novel_id
	comments map[string][]*Comment // keyed by chapter_id
	library  map[string][]*LibraryEntry // keyed by user_id
	progress map[string]*ReadingProgress // keyed by "user_id:novel_id"
	genres   []string

	userByEmail map[string]string
	novelBySlug map[string]string
}

func NewStore() *Store {
	s := &Store{
		users:       make(map[string]*User),
		novels:      make(map[string]*Novel),
		chapters:    make(map[string][]*Chapter),
		comments:    make(map[string][]*Comment),
		library:     make(map[string][]*LibraryEntry),
		progress:    make(map[string]*ReadingProgress),
		userByEmail: make(map[string]string),
		novelBySlug: make(map[string]string),
	}
	s.seed()
	return s
}

func genID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func hashPassword(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h)
}


func (s *Store) CreateUser(username, email, password, role string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	email = strings.ToLower(email)
	if _, ok := s.userByEmail[email]; ok {
		return nil, fmt.Errorf("email already exists")
	}
	for _, u := range s.users {
		if strings.EqualFold(u.Username, username) {
			return nil, fmt.Errorf("username already exists")
		}
	}
	u := &User{
		ID: genID(), Username: username, Email: email,
		PasswordHash: hashPassword(password), Role: role,
		CreatedAt: time.Now(),
	}
	s.users[u.ID] = u
	s.userByEmail[email] = u.ID
	return u, nil
}

func (s *Store) AuthenticateUser(email, password string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uid, ok := s.userByEmail[strings.ToLower(email)]
	if !ok {
		return nil, fmt.Errorf("invalid credentials")
	}
	u := s.users[uid]
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return u, nil
}

func (s *Store) GetUser(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) ListUsers() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*User
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

func (s *Store) UpdateUserRole(id, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}
	u.Role = role
	return nil
}

func (s *Store) CreateNovel(n *Novel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.novelBySlug[n.Slug]; ok {
		return fmt.Errorf("novel slug already exists")
	}
	n.ID = genID()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	s.novels[n.ID] = n
	s.novelBySlug[n.Slug] = n.ID
	return nil
}

func (s *Store) GetNovelBySlug(slug string) (*Novel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.novelBySlug[slug]
	if !ok {
		return nil, fmt.Errorf("novel not found")
	}
	return s.novels[id], nil
}

func (s *Store) GetNovelByID(id string) (*Novel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.novels[id]
	if !ok {
		return nil, fmt.Errorf("novel not found")
	}
	return n, nil
}

func (s *Store) ListNovels(genre, status, sortBy string, page, pageSize int) ([]*Novel, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var all []*Novel
	for _, n := range s.novels {
		if genre != "" {
			found := false
			for _, g := range n.Genres {
				if strings.EqualFold(g, genre) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if status != "" && n.Status != status {
			continue
		}
		all = append(all, n)
	}
	sort.Slice(all, func(i, j int) bool {
		switch sortBy {
		case "title":
			return all[i].Title < all[j].Title
		case "chapters":
			return all[i].TotalChapters > all[j].TotalChapters
		default:
			return all[i].UpdatedAt.After(all[j].UpdatedAt)
		}
	})
	total := len(all)
	start := (page - 1) * pageSize
	if start >= total {
		return nil, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return all[start:end], total
}

func (s *Store) UpdateNovel(id string, title, synopsis, author, status, coverURL string, genres []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.novels[id]
	if !ok {
		return fmt.Errorf("novel not found")
	}
	if title != "" {
		n.Title = title
	}
	if synopsis != "" {
		n.Synopsis = synopsis
	}
	if author != "" {
		n.Author = author
	}
	if status != "" {
		n.Status = status
	}
	if coverURL != "" {
		n.CoverURL = coverURL
	}
	if genres != nil {
		n.Genres = genres
	}
	n.UpdatedAt = time.Now()
	return nil
}

func (s *Store) DeleteNovel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.novels[id]
	if !ok {
		return fmt.Errorf("novel not found")
	}
	delete(s.novelBySlug, n.Slug)
	delete(s.novels, id)
	delete(s.chapters, id)
	return nil
}

func (s *Store) CreateChapter(ch *Chapter) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch.ID = genID()
	ch.WordCount = len(strings.Fields(ch.Content))
	ch.CreatedAt = time.Now()
	s.chapters[ch.NovelID] = append(s.chapters[ch.NovelID], ch)
	if n, ok := s.novels[ch.NovelID]; ok {
		n.TotalChapters = len(s.chapters[ch.NovelID])
		n.UpdatedAt = time.Now()
	}
	return nil
}

func (s *Store) GetChapter(novelSlug string, number int) (*Chapter, *Novel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nid, ok := s.novelBySlug[novelSlug]
	if !ok {
		return nil, nil, fmt.Errorf("novel not found")
	}
	novel := s.novels[nid]
	for _, ch := range s.chapters[nid] {
		if ch.Number == number {
			return ch, novel, nil
		}
	}
	return nil, nil, fmt.Errorf("chapter not found")
}

func (s *Store) ListChapters(novelSlug string) ([]*Chapter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	nid, ok := s.novelBySlug[novelSlug]
	if !ok {
		return nil, fmt.Errorf("novel not found")
	}
	chs := s.chapters[nid]
	sort.Slice(chs, func(i, j int) bool { return chs[i].Number < chs[j].Number })
	return chs, nil
}

func (s *Store) ListAllChapters() []*struct {
	Chapter    *Chapter
	NovelTitle string
} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*struct {
		Chapter    *Chapter
		NovelTitle string
	}
	for nid, chs := range s.chapters {
		title := ""
		if n, ok := s.novels[nid]; ok {
			title = n.Title
		}
		for _, ch := range chs {
			out = append(out, &struct {
				Chapter    *Chapter
				NovelTitle string
			}{Chapter: ch, NovelTitle: title})
		}
	}
	return out
}

func (s *Store) UpdateChapter(id, title, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, chs := range s.chapters {
		for _, ch := range chs {
			if ch.ID == id {
				if title != "" {
					ch.Title = title
				}
				if content != "" {
					ch.Content = content
					ch.WordCount = len(strings.Fields(content))
				}
				return nil
			}
		}
	}
	return fmt.Errorf("chapter not found")
}

func (s *Store) DeleteChapter(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for nid, chs := range s.chapters {
		for i, ch := range chs {
			if ch.ID == id {
				s.chapters[nid] = append(chs[:i], chs[i+1:]...)
				if n, ok := s.novels[nid]; ok {
					n.TotalChapters = len(s.chapters[nid])
				}
				return nil
			}
		}
	}
	return fmt.Errorf("chapter not found")
}

func (s *Store) CreateComment(c *Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = genID()
	c.CreatedAt = time.Now()
	s.comments[c.ChapterID] = append(s.comments[c.ChapterID], c)
}

func (s *Store) ListComments(chapterID string) []*Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.comments[chapterID]
}

func (s *Store) LikeComment(commentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, cms := range s.comments {
		for _, c := range cms {
			if c.ID == commentID {
				c.Likes++
				return
			}
		}
	}
}


func (s *Store) AddToLibrary(userID, novelID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.novels[novelID]
	if !ok {
		return fmt.Errorf("novel not found")
	}
	for _, e := range s.library[userID] {
		if e.NovelID == novelID {
			return fmt.Errorf("already in library")
		}
	}
	s.library[userID] = append(s.library[userID], &LibraryEntry{
		UserID: userID, NovelID: novelID, NovelTitle: n.Title,
		NovelSlug: n.Slug, Status: "reading", Total: n.TotalChapters,
	})
	return nil
}

func (s *Store) GetLibrary(userID string) []*LibraryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.library[userID]
}

func (s *Store) UpdateLibraryStatus(userID, novelID, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.library[userID] {
		if e.NovelID == novelID {
			e.Status = status
			return
		}
	}
}

func (s *Store) SaveProgress(userID, novelID string, chapterNum, scroll int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := userID + ":" + novelID
	s.progress[key] = &ReadingProgress{
		UserID: userID, NovelID: novelID,
		ChapterNumber: chapterNum, ScrollPosition: scroll,
	}
}

func (s *Store) GetProgress(userID, novelID string) *ReadingProgress {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.progress[userID+":"+novelID]
}

func (s *Store) SearchNovels(q string) []*Novel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q = strings.ToLower(q)
	var results []*Novel
	for _, n := range s.novels {
		if strings.Contains(strings.ToLower(n.Title), q) ||
			strings.Contains(strings.ToLower(n.Author), q) ||
			strings.Contains(strings.ToLower(n.Synopsis), q) {
			results = append(results, n)
		}
	}
	return results
}

func (s *Store) GetGenres() []string {
	return s.genres
}
