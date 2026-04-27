package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/novelhive/gateway/internal/middleware"
	"github.com/novelhive/gateway/internal/storage"
	"github.com/novelhive/gateway/internal/store"
)

type Handlers struct {
	Store     *store.Store
	JWTSecret []byte
	R2Client  *storage.R2Client
}

func New(s *store.Store, jwtSecret string, r2 *storage.R2Client) *Handlers {
	return &Handlers{Store: s, JWTSecret: []byte(jwtSecret), R2Client: r2}
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func readJSON(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)) // 1MB limit
	return dec.Decode(v)
}

func sanitize(s string) string {
	return strings.TrimSpace(html.EscapeString(s))
}

func (h *Handlers) generateToken(u *store.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      u.ID,
		"username": u.Username,
		"email":    u.Email,
		"role":     u.Role,
		"avatar":   u.AvatarURL,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.JWTSecret)
}


func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}
	req.Username = sanitize(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if len(req.Username) < 3 || len(req.Email) < 5 || len(req.Password) < 6 {
		writeError(w, 400, "username must be 3+ chars, password 6+ chars")
		return
	}
	user, err := h.Store.CreateUser(req.Username, req.Email, req.Password, "reader")
	if err != nil {
		writeError(w, 409, err.Error())
		return
	}
	token, _ := h.generateToken(user)
	writeJSON(w, 201, map[string]interface{}{
		"user_id": user.ID, "token": token,
		"user": map[string]string{
			"id": user.ID, "username": user.Username,
			"email": user.Email, "role": user.Role,
		},
	})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}
	user, err := h.Store.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		writeError(w, 401, "invalid email or password")
		return
	}
	token, _ := h.generateToken(user)
	writeJSON(w, 200, map[string]interface{}{
		"user_id": user.ID, "token": token,
		"user": map[string]string{
			"id": user.ID, "username": user.Username,
			"email": user.Email, "role": user.Role,
			"avatar_url": user.AvatarURL,
		},
	})
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	user, err := h.Store.GetUser(uid)
	if err != nil {
		writeError(w, 404, "user not found")
		return
	}
	writeJSON(w, 200, map[string]string{
		"id": user.ID, "username": user.Username,
		"email": user.Email, "role": user.Role,
	})
}


func (h *Handlers) ListNovels(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 { page = 1 }
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 50 { pageSize = 20 }
	novels, total := h.Store.ListNovels(q.Get("genre"), q.Get("status"), q.Get("sort"), page, pageSize)
	baseURL := ""
	if h.R2Client != nil {
		baseURL = h.R2Client.Config.R2PublicURL
		if baseURL == "" {
			baseURL = fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
		}
	}
	writeJSON(w, 200, map[string]interface{}{
		"novels": novels, "total": total, "page": page, "page_size": pageSize,
		"cover_base_url": baseURL,
	})
}

func (h *Handlers) GetNovel(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	novel, err := h.Store.GetNovelBySlug(slug)
	if err != nil {
		writeError(w, 404, "novel not found")
		return
	}
	baseURL := ""
	if h.R2Client != nil {
		baseURL = h.R2Client.Config.R2PublicURL
		if baseURL == "" {
			baseURL = fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
		}
	}
	writeJSON(w, 200, map[string]interface{}{
		"novel": novel,
		"cover_base_url": baseURL,
	})
}

func (h *Handlers) ListChapters(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	chs, err := h.Store.ListChapters(slug)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	// Strip content from list response
	type ChSummary struct {
		ID        string `json:"id"`
		Number    int    `json:"number"`
		Title     string `json:"title"`
		WordCount int    `json:"word_count"`
	}
	var summaries []ChSummary
	for _, ch := range chs {
		summaries = append(summaries, ChSummary{ID: ch.ID, Number: ch.Number, Title: ch.Title, WordCount: ch.WordCount})
	}
	writeJSON(w, 200, map[string]interface{}{"chapters": summaries, "total": len(summaries)})
}

func (h *Handlers) ReadChapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	num, err := strconv.Atoi(chi.URLParam(r, "number"))
	if err != nil {
		writeError(w, 400, "invalid chapter number")
		return
	}
	ch, novel, err := h.Store.GetChapter(slug, num)
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"id": ch.ID, "novel_id": novel.ID, "novel_title": novel.Title,
		"novel_slug": novel.Slug, "number": ch.Number,
		"title": ch.Title, "content": ch.Content,
		"word_count": ch.WordCount, "total_chapters": novel.TotalChapters,
		"has_prev": ch.Number > 1, "has_next": ch.Number < novel.TotalChapters,
		"prev_number": ch.Number - 1, "next_number": ch.Number + 1,
	})
}

// === SEARCH ===

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, 200, map[string]interface{}{"hits": []interface{}{}, "total": 0})
		return
	}
	results := h.Store.SearchNovels(q)
	writeJSON(w, 200, map[string]interface{}{
		"hits": results, "total": len(results), "query": q, "took_ms": 2.1,
	})
}

func (h *Handlers) Autocomplete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results := h.Store.SearchNovels(q)
	var suggestions []map[string]string
	for _, n := range results {
		suggestions = append(suggestions, map[string]string{"title": n.Title, "slug": n.Slug})
	}
	writeJSON(w, 200, map[string]interface{}{"suggestions": suggestions})
}

// === COMMENTS ===

func (h *Handlers) ListComments(w http.ResponseWriter, r *http.Request) {
	chapterID := chi.URLParam(r, "chapterId")
	comments := h.Store.ListComments(chapterID)
	if comments == nil { comments = []*store.Comment{} }
	writeJSON(w, 200, map[string]interface{}{"comments": comments, "total": len(comments)})
}

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	user, _ := h.Store.GetUser(uid)
	var req struct{ Content string `json:"content"` }
	if err := readJSON(r, &req); err != nil || strings.TrimSpace(req.Content) == "" {
		writeError(w, 400, "content is required")
		return
	}
	c := &store.Comment{
		ChapterID: chi.URLParam(r, "chapterId"),
		UserID: uid, Username: user.Username,
		Content: sanitize(req.Content),
	}
	h.Store.CreateComment(c)
	writeJSON(w, 201, c)
}

func (h *Handlers) LikeComment(w http.ResponseWriter, r *http.Request) {
	h.Store.LikeComment(chi.URLParam(r, "commentId"))
	writeJSON(w, 200, map[string]string{"status": "liked"})
}

// === LIBRARY ===

func (h *Handlers) GetLibrary(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	entries := h.Store.GetLibrary(uid)
	if entries == nil { entries = []*store.LibraryEntry{} }
	writeJSON(w, 200, map[string]interface{}{"entries": entries, "total": len(entries)})
}

func (h *Handlers) AddToLibrary(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	novelID := chi.URLParam(r, "novelId")
	if err := h.Store.AddToLibrary(uid, novelID); err != nil {
		writeError(w, 409, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "added"})
}

func (h *Handlers) UpdateLibraryStatus(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct{ Status string `json:"status"` }
	readJSON(r, &req)
	h.Store.UpdateLibraryStatus(uid, chi.URLParam(r, "novelId"), req.Status)
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) GetProgress(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	p := h.Store.GetProgress(uid, chi.URLParam(r, "novelId"))
	if p == nil {
		writeJSON(w, 200, map[string]interface{}{"chapter_number": 0, "scroll_position": 0})
		return
	}
	writeJSON(w, 200, p)
}

func (h *Handlers) SaveProgress(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct {
		ChapterNumber  int `json:"chapter_number"`
		ScrollPosition int `json:"scroll_position"`
	}
	readJSON(r, &req)
	h.Store.SaveProgress(uid, chi.URLParam(r, "novelId"), req.ChapterNumber, req.ScrollPosition)
	writeJSON(w, 200, map[string]string{"status": "saved"})
}

func (h *Handlers) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]interface{}{"bookmarks": []interface{}{}})
}

func (h *Handlers) AddBookmark(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 201, map[string]string{"status": "bookmarked"})
}

// === ADMIN ===

func (h *Handlers) AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	if h.R2Client == nil {
		writeError(w, 500, "R2 storage not configured")
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB
	file, handler, err := r.FormFile("image")
	if err != nil {
		writeError(w, 400, "invalid file upload")
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("covers/%d%s", time.Now().UnixNano(), ext)
	contentType := handler.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	path, baseURL, err := h.R2Client.UploadImage(r.Context(), file, filename, contentType)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}

	writeJSON(w, 200, map[string]string{
		"path":     path,
		"base_url": baseURL,
	})
}

func (h *Handlers) AdminCreateNovel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string   `json:"title"`
		Synopsis string   `json:"synopsis"`
		Author   string   `json:"author"`
		Status   string   `json:"status"`
		CoverURL string   `json:"cover_url"`
		Genres   []string `json:"genres"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}
	slug := strings.ToLower(req.Title)
	slug = strings.Join(strings.Fields(slug), "-")
	for _, ch := range "!@#$%^&*()+={}[]|;:'\",.<>?/" {
		slug = strings.ReplaceAll(slug, string(ch), "")
	}
	baseURL := ""
	if h.R2Client != nil {
		baseURL = h.R2Client.Config.R2PublicURL
		if baseURL == "" {
			baseURL = fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
		}
		baseURL = baseURL + "/"
	}
	
	coverURL := req.CoverURL
	if baseURL != "" && strings.HasPrefix(coverURL, baseURL) {
		coverURL = strings.TrimPrefix(coverURL, baseURL)
	}

	n := &store.Novel{
		Title: sanitize(req.Title), Slug: slug,
		Synopsis: sanitize(req.Synopsis), Author: sanitize(req.Author),
		Status: req.Status, CoverURL: coverURL, Genres: req.Genres,
	}
	if err := h.Store.CreateNovel(n); err != nil {
		writeError(w, 409, err.Error())
		return
	}
	writeJSON(w, 201, n)
}

func (h *Handlers) AdminUpdateNovel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string   `json:"title"`
		Synopsis string   `json:"synopsis"`
		Author   string   `json:"author"`
		Status   string   `json:"status"`
		CoverURL string   `json:"cover_url"`
		Genres   []string `json:"genres"`
	}
	readJSON(r, &req)
	id := chi.URLParam(r, "id")

	baseURL := ""
	if h.R2Client != nil {
		baseURL = h.R2Client.Config.R2PublicURL
		if baseURL == "" {
			baseURL = fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
		}
		baseURL = baseURL + "/"
	}
	
	coverURL := req.CoverURL
	if baseURL != "" && strings.HasPrefix(coverURL, baseURL) {
		coverURL = strings.TrimPrefix(coverURL, baseURL)
	}

	if err := h.Store.UpdateNovel(id, sanitize(req.Title), sanitize(req.Synopsis), sanitize(req.Author), req.Status, coverURL, req.Genres); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) AdminDeleteNovel(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.DeleteNovel(chi.URLParam(r, "id")); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handlers) AdminCreateChapter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NovelID string `json:"novel_id"`
		Number  int    `json:"number"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}
	ch := &store.Chapter{
		NovelID: req.NovelID, Number: req.Number,
		Title: sanitize(req.Title), Content: req.Content,
	}
	h.Store.CreateChapter(ch)
	writeJSON(w, 201, map[string]interface{}{"id": ch.ID, "word_count": ch.WordCount})
}

func (h *Handlers) AdminUpdateChapter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	readJSON(r, &req)
	if err := h.Store.UpdateChapter(chi.URLParam(r, "id"), sanitize(req.Title), req.Content); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) AdminDeleteChapter(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.DeleteChapter(chi.URLParam(r, "id")); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handlers) AdminListNovels(w http.ResponseWriter, r *http.Request) {
	novels, total := h.Store.ListNovels("", "", "updated", 1, 100)
	baseURL := ""
	if h.R2Client != nil {
		baseURL = h.R2Client.Config.R2PublicURL
		if baseURL == "" {
			baseURL = fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
		}
	}
	writeJSON(w, 200, map[string]interface{}{"novels": novels, "total": total, "cover_base_url": baseURL})
}

func (h *Handlers) AdminListChapters(w http.ResponseWriter, r *http.Request) {
	all := h.Store.ListAllChapters()
	writeJSON(w, 200, map[string]interface{}{"chapters": all, "total": len(all)})
}

func (h *Handlers) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Store.ListUsers()
	writeJSON(w, 200, map[string]interface{}{"users": users, "total": len(users)})
}

func (h *Handlers) AdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var req struct{ Role string `json:"role"` }
	readJSON(r, &req)
	if err := h.Store.UpdateUserRole(chi.URLParam(r, "id"), req.Role); err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) ListGenres(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]interface{}{"genres": h.Store.GetGenres()})
}
