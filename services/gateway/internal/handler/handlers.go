package handler

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/novelhive/gateway/internal/clients"
	"github.com/novelhive/gateway/internal/middleware"
	"github.com/novelhive/gateway/internal/storage"
	commentv1 "github.com/novelhive/proto/comment/v1"
	libraryv1 "github.com/novelhive/proto/library/v1"
	notificationv1 "github.com/novelhive/proto/notification/v1"
	novelv1 "github.com/novelhive/proto/novel/v1"
	userv1 "github.com/novelhive/proto/user/v1"
)
...
func (h *Handlers) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	resp, err := h.Clients.Notification.GetNotifications(r.Context(), &notificationv1.GetNotificationsRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		writeError(w, 500, "failed to get notifications")
		return
	}

	writeJSON(w, 200, resp)
}

func (h *Handlers) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	resp, err := h.Clients.Notification.MarkAsRead(r.Context(), &notificationv1.MarkAsReadRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		writeError(w, 500, "failed to mark notification as read")
		return
	}

	writeJSON(w, 200, resp)
}

func (h *Handlers) RegisterFCMToken(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req struct {
		Token string `json:"token"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	resp, err := h.Clients.Notification.RegisterFCMToken(r.Context(), &notificationv1.RegisterFCMTokenRequest{
		UserId: userID,
		Token:  req.Token,
	})
	if err != nil {
		writeError(w, 500, "failed to register FCM token")
		return
	}

	writeJSON(w, 200, resp)
}

type Handlers struct {
	Clients   *clients.Clients
	JWTSecret []byte
	R2Client  *storage.R2Client
}

func New(c *clients.Clients, jwtSecret string, r2 *storage.R2Client) *Handlers {
	return &Handlers{Clients: c, JWTSecret: []byte(jwtSecret), R2Client: r2}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func readJSON(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	return dec.Decode(v)
}

func sanitize(s string) string {
	return strings.TrimSpace(html.EscapeString(s))
}

func (h *Handlers) r2BaseURL() string {
	if h.R2Client == nil {
		return ""
	}
	if h.R2Client.Config.R2PublicURL != "" {
		return h.R2Client.Config.R2PublicURL
	}
	return fmt.Sprintf("%s/%s", h.R2Client.Config.R2Endpoint, h.R2Client.Config.R2BucketName)
}

func (h *Handlers) generateAccessToken(userID, username, email, role, avatarURL string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"email":    email,
		"role":     role,
		"avatar":   avatarURL,
		"type":     "access",
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
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
		writeError(w, 400, "username ≥3 chars, password ≥6 chars")
		return
	}

	resp, err := h.Clients.User.Register(r.Context(), &userv1.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, 409, err.Error())
		return
	}
	writeJSON(w, 201, map[string]interface{}{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
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
	resp, err := h.Clients.User.Login(r.Context(), &userv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, 401, "invalid email or password")
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"user_id":       resp.UserId,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"user": map[string]string{
			"id":         resp.Profile.Id,
			"username":   resp.Profile.Username,
			"email":      resp.Profile.Email,
			"role":       resp.Profile.Role,
			"avatar_url": resp.Profile.AvatarUrl,
		},
	})
}

func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := readJSON(r, &req); err != nil || req.RefreshToken == "" {
		writeError(w, 400, "refresh_token is required")
		return
	}

	resp, err := h.Clients.User.RefreshToken(r.Context(), &userv1.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		writeError(w, 401, "invalid or expired refresh token")
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	resp, err := h.Clients.User.GetProfile(r.Context(), &userv1.GetProfileRequest{UserId: uid})
	if err != nil {
		writeError(w, 404, "user not found")
		return
	}
	writeJSON(w, 200, map[string]string{
		"id":         resp.Id,
		"username":   resp.Username,
		"email":      resp.Email,
		"role":       resp.Role,
		"avatar_url": resp.AvatarUrl,
	})
}

func (h *Handlers) ListNovels(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	resp, err := h.Clients.Novel.ListNovels(r.Context(), &novelv1.ListNovelsRequest{
		Page:      int32(page),
		PageSize:  int32(pageSize),
		GenreSlug: q.Get("genre"),
		Status:    q.Get("status"),
		SortBy:    q.Get("sort"),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"novels":         resp.Novels,
		"total":          resp.Total,
		"page":           page,
		"page_size":      pageSize,
		"cover_base_url": h.r2BaseURL(),
	})
}

func (h *Handlers) GetNovel(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	resp, err := h.Clients.Novel.GetNovel(r.Context(), &novelv1.GetNovelRequest{Slug: slug})
	if err != nil {
		writeError(w, 404, "novel not found")
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"novel":          resp,
		"cover_base_url": h.r2BaseURL(),
	})
}

func (h *Handlers) ListChapters(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	novel, err := h.Clients.Novel.GetNovel(r.Context(), &novelv1.GetNovelRequest{Slug: slug})
	if err != nil {
		writeError(w, 404, "novel not found")
		return
	}
	resp, err := h.Clients.Novel.ListChapters(r.Context(), &novelv1.ListChaptersRequest{
		NovelId:  novel.Id,
		Page:     1,
		PageSize: 500,
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"chapters": resp.Chapters,
		"total":    resp.Total,
	})
}

func (h *Handlers) ReadChapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	num, err := strconv.Atoi(chi.URLParam(r, "number"))
	if err != nil {
		writeError(w, 400, "invalid chapter number")
		return
	}
	ch, err := h.Clients.Novel.GetChapter(r.Context(), &novelv1.GetChapterRequest{
		NovelSlug:     slug,
		ChapterNumber: int32(num),
	})
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}
	novel, _ := h.Clients.Novel.GetNovel(r.Context(), &novelv1.GetNovelRequest{Slug: slug})
	var totalChapters int32
	if novel != nil {
		totalChapters = novel.TotalChapters
	}
	writeJSON(w, 200, map[string]interface{}{
		"id":       ch.Id,
		"novel_id": ch.NovelId,
		"novel_title": func() string {
			if novel != nil {
				return novel.Title
			}
			return ""
		}(),
		"novel_slug":     slug,
		"number":         ch.Number,
		"title":          ch.Title,
		"content":        ch.Content,
		"word_count":     ch.WordCount,
		"total_chapters": totalChapters,
		"has_prev":       ch.Number > 1,
		"has_next":       ch.Number < totalChapters,
		"prev_number":    ch.Number - 1,
		"next_number":    ch.Number + 1,
	})
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, 200, map[string]interface{}{"hits": []interface{}{}, "total": 0})
		return
	}
	// Fallback: list novels and filter in-process (search-service not wired yet)
	resp, err := h.Clients.Novel.ListNovels(r.Context(), &novelv1.ListNovelsRequest{Page: 1, PageSize: 50})
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{"hits": []interface{}{}, "total": 0})
		return
	}
	ql := strings.ToLower(q)
	var hits []*novelv1.Novel
	for _, n := range resp.Novels {
		if strings.Contains(strings.ToLower(n.Title), ql) ||
			strings.Contains(strings.ToLower(n.Author), ql) ||
			strings.Contains(strings.ToLower(n.Synopsis), ql) {
			hits = append(hits, n)
		}
	}
	writeJSON(w, 200, map[string]interface{}{
		"hits": hits, "total": len(hits), "query": q, "took_ms": 1.0,
	})
}

func (h *Handlers) Autocomplete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	resp, _ := h.Clients.Novel.ListNovels(r.Context(), &novelv1.ListNovelsRequest{Page: 1, PageSize: 50})
	ql := strings.ToLower(q)
	var suggestions []map[string]string
	if resp != nil {
		for _, n := range resp.Novels {
			if strings.Contains(strings.ToLower(n.Title), ql) {
				suggestions = append(suggestions, map[string]string{"title": n.Title, "slug": n.Slug})
			}
		}
	}
	writeJSON(w, 200, map[string]interface{}{"suggestions": suggestions})
}


func (h *Handlers) ListGenres(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Clients.Novel.ListGenres(r.Context(), &novelv1.ListGenresRequest{})
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{"genres": []interface{}{}})
		return
	}
	var genres []string
	for _, g := range resp.Genres {
		genres = append(genres, g.Name)
	}
	writeJSON(w, 200, map[string]interface{}{"genres": genres})
}


func (h *Handlers) ListComments(w http.ResponseWriter, r *http.Request) {
	chapterID := chi.URLParam(r, "chapterId")
	resp, err := h.Clients.Comment.ListComments(r.Context(), &commentv1.ListCommentsRequest{
		ChapterId: chapterID,
		Page:      1,
		PageSize:  50,
	})
	if err != nil || resp == nil {
		writeJSON(w, 200, map[string]interface{}{"comments": []interface{}{}, "total": 0})
		return
	}
	writeJSON(w, 200, map[string]interface{}{"comments": resp.Comments, "total": resp.Total})
}

func (h *Handlers) CreateComment(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct {
		Content string `json:"content"`
	}
	if err := readJSON(r, &req); err != nil || strings.TrimSpace(req.Content) == "" {
		writeError(w, 400, "content is required")
		return
	}
	resp, err := h.Clients.Comment.CreateComment(r.Context(), &commentv1.CreateCommentRequest{
		ChapterId: chi.URLParam(r, "chapterId"),
		UserId:    uid,
		Content:   sanitize(req.Content),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, resp)
}

func (h *Handlers) LikeComment(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	resp, _ := h.Clients.Comment.LikeComment(r.Context(), &commentv1.LikeCommentRequest{
		CommentId: chi.URLParam(r, "commentId"),
		UserId:    uid,
	})
	writeJSON(w, 200, resp)
}

func (h *Handlers) GetLibrary(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())

	const (
		libraryPageSize = 50
		libraryMaxPages = 20
	)

	var total int32
	allEntries := make([]*libraryv1.LibraryEntry, 0, libraryPageSize)
	for page := 1; page <= libraryMaxPages; page++ {
		resp, err := h.Clients.Library.GetLibrary(r.Context(), &libraryv1.GetLibraryRequest{
			UserId:   uid,
			Page:     int32(page),
			PageSize: libraryPageSize,
		})
		if err != nil || resp == nil {
			break
		}
		if page == 1 {
			total = resp.Total
		}
		if len(resp.Entries) == 0 {
			break
		}
		allEntries = append(allEntries, resp.Entries...)
		if total > 0 && int32(len(allEntries)) >= total {
			break
		}
		if len(resp.Entries) < libraryPageSize {
			break
		}
	}

	if len(allEntries) == 0 {
		writeJSON(w, 200, map[string]interface{}{"entries": []interface{}{}, "total": 0, "cover_base_url": h.r2BaseURL()})
		return
	}

	ids := make([]string, 0, len(allEntries))
	seen := make(map[string]struct{}, len(allEntries))
	for _, e := range allEntries {
		if e == nil || e.NovelId == "" {
			continue
		}
		if _, ok := seen[e.NovelId]; ok {
			continue
		}
		seen[e.NovelId] = struct{}{}
		ids = append(ids, e.NovelId)
	}

	novelsByID := make(map[string]*novelv1.Novel, len(ids))
	const novelBatchSize = 200
	for start := 0; start < len(ids); start += novelBatchSize {
		end := start + novelBatchSize
		if end > len(ids) {
			end = len(ids)
		}
		novelsResp, nerr := h.Clients.Novel.GetNovelsByIds(r.Context(), &novelv1.GetNovelsByIdsRequest{Ids: ids[start:end]})
		if nerr != nil || novelsResp == nil {
			continue
		}
		for _, n := range novelsResp.Novels {
			if n != nil && n.Id != "" {
				novelsByID[n.Id] = n
			}
		}
	}

	const progressConcurrency = 8
	progressByNovelID := make(map[string]int32, len(ids))
	if len(ids) > 0 {
		sem := make(chan struct{}, progressConcurrency)
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, id := range ids {
			id := id
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				p, perr := h.Clients.Library.GetProgress(r.Context(), &libraryv1.GetProgressRequest{UserId: uid, NovelId: id})
				if perr != nil || p == nil {
					return
				}
				mu.Lock()
				progressByNovelID[id] = p.ChapterNumber
				mu.Unlock()
			}()
		}
		wg.Wait()
	}

	entries := make([]map[string]interface{}, 0, len(allEntries))
	for _, e := range allEntries {
		if e == nil {
			continue
		}
		entry := map[string]interface{}{
			"id":         e.Id,
			"user_id":    e.UserId,
			"novel_id":   e.NovelId,
			"status":     e.Status,
			"created_at": e.CreatedAt,
		}
		if n := novelsByID[e.NovelId]; n != nil {
			entry["novel_title"] = n.Title
			entry["novel_slug"] = n.Slug
			entry["total"] = n.TotalChapters
			entry["cover_url"] = n.CoverUrl
			entry["author"] = n.Author
		}
		if p, ok := progressByNovelID[e.NovelId]; ok {
			entry["progress"] = p
		} else {
			entry["progress"] = int32(0)
		}
		entries = append(entries, entry)
	}

	if total == 0 {
		total = int32(len(entries))
	}

	writeJSON(w, 200, map[string]interface{}{
		"entries":        entries,
		"total":          total,
		"cover_base_url": h.r2BaseURL(),
	})
}

func (h *Handlers) GetLibraryEntry(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	novelID := chi.URLParam(r, "novelId")
	if novelID == "" {
		writeError(w, 400, "novelId is required")
		return
	}

	const (
		pageSize = 50
		maxPages = 20
	)

	var entry *libraryv1.LibraryEntry
	for page := 1; page <= maxPages; page++ {
		resp, err := h.Clients.Library.GetLibrary(r.Context(), &libraryv1.GetLibraryRequest{UserId: uid, Page: int32(page), PageSize: pageSize})
		if err != nil || resp == nil || len(resp.Entries) == 0 {
			break
		}
		for _, e := range resp.Entries {
			if e != nil && e.NovelId == novelID {
				entry = e
				break
			}
		}
		if entry != nil {
			break
		}
		if len(resp.Entries) < pageSize {
			break
		}
	}

	if entry == nil {
		writeJSON(w, 200, map[string]interface{}{"in_library": false, "cover_base_url": h.r2BaseURL()})
		return
	}

	var novel map[string]interface{}
	if novelsResp, err := h.Clients.Novel.GetNovelsByIds(r.Context(), &novelv1.GetNovelsByIdsRequest{Ids: []string{novelID}}); err == nil && novelsResp != nil {
		if len(novelsResp.Novels) > 0 && novelsResp.Novels[0] != nil {
			n := novelsResp.Novels[0]
			novel = map[string]interface{}{
				"id":             n.Id,
				"title":          n.Title,
				"slug":           n.Slug,
				"cover_url":      n.CoverUrl,
				"author":         n.Author,
				"total_chapters": n.TotalChapters,
			}
		}
	}

	progress := int32(0)
	if p, err := h.Clients.Library.GetProgress(r.Context(), &libraryv1.GetProgressRequest{UserId: uid, NovelId: novelID}); err == nil && p != nil {
		progress = p.ChapterNumber
	}

	entryJSON := map[string]interface{}{
		"id":         entry.Id,
		"user_id":    entry.UserId,
		"novel_id":   entry.NovelId,
		"status":     entry.Status,
		"created_at": entry.CreatedAt,
	}

	writeJSON(w, 200, map[string]interface{}{
		"in_library":     true,
		"entry":          entryJSON,
		"novel":          novel,
		"progress":       progress,
		"cover_base_url": h.r2BaseURL(),
	})
}

func (h *Handlers) AddToLibrary(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	_, err := h.Clients.Library.AddToLibrary(r.Context(), &libraryv1.AddToLibraryRequest{
		UserId:  uid,
		NovelId: chi.URLParam(r, "novelId"),
	})
	if err != nil {
		writeError(w, 409, err.Error())
		return
	}
	writeJSON(w, 201, map[string]string{"status": "added"})
}

func (h *Handlers) RemoveFromLibrary(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	_, err := h.Clients.Library.RemoveFromLibrary(r.Context(), &libraryv1.RemoveFromLibraryRequest{
		UserId:  uid,
		NovelId: chi.URLParam(r, "novelId"),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "removed"})
}

func (h *Handlers) UpdateLibraryStatus(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct {
		Status string `json:"status"`
	}
	readJSON(r, &req)
	h.Clients.Library.UpdateStatus(r.Context(), &libraryv1.UpdateStatusRequest{
		UserId:  uid,
		NovelId: chi.URLParam(r, "novelId"),
		Status:  req.Status,
	})
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) GetProgress(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	resp, err := h.Clients.Library.GetProgress(r.Context(), &libraryv1.GetProgressRequest{
		UserId:  uid,
		NovelId: chi.URLParam(r, "novelId"),
	})
	if err != nil || resp == nil {
		writeJSON(w, 200, map[string]interface{}{"chapter_number": 0, "scroll_position": 0})
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"chapter_number":  resp.ChapterNumber,
		"scroll_position": resp.ScrollPosition,
	})
}

func (h *Handlers) SaveProgress(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct {
		ChapterNumber  int32   `json:"chapter_number"`
		ScrollPosition float32 `json:"scroll_position"`
	}
	readJSON(r, &req)
	h.Clients.Library.SaveProgress(r.Context(), &libraryv1.SaveProgressRequest{
		UserId:         uid,
		NovelId:        chi.URLParam(r, "novelId"),
		ChapterNumber:  req.ChapterNumber,
		ScrollPosition: req.ScrollPosition,
	})
	writeJSON(w, 200, map[string]string{"status": "saved"})
}

func (h *Handlers) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	novelID := r.URL.Query().Get("novel_id")
	resp, err := h.Clients.Library.GetBookmarks(r.Context(), &libraryv1.GetBookmarksRequest{
		UserId:  uid,
		NovelId: novelID,
	})
	if err != nil || resp == nil {
		writeJSON(w, 200, map[string]interface{}{"bookmarks": []interface{}{}})
		return
	}
	writeJSON(w, 200, map[string]interface{}{"bookmarks": resp.Bookmarks})
}

func (h *Handlers) AddBookmark(w http.ResponseWriter, r *http.Request) {
	uid := middleware.GetUserID(r.Context())
	var req struct {
		NovelID   string `json:"novel_id"`
		ChapterID string `json:"chapter_id"`
		Note      string `json:"note"`
	}
	readJSON(r, &req)
	resp, err := h.Clients.Library.AddBookmark(r.Context(), &libraryv1.AddBookmarkRequest{
		UserId:    uid,
		NovelId:   req.NovelID,
		ChapterId: req.ChapterID,
		Note:      req.Note,
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, resp)
}

func (h *Handlers) AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	if h.R2Client == nil {
		writeError(w, 500, "R2 storage not configured")
		return
	}
	r.ParseMultipartForm(10 << 20)
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
	writeJSON(w, 200, map[string]string{"path": path, "base_url": baseURL})
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

	baseURL := h.r2BaseURL()
	coverURL := req.CoverURL
	if baseURL != "" && strings.HasPrefix(coverURL, baseURL+"/") {
		coverURL = strings.TrimPrefix(coverURL, baseURL+"/")
	}

	var genreIDs []int32
	genreResp, _ := h.Clients.Novel.ListGenres(r.Context(), &novelv1.ListGenresRequest{})
	if genreResp != nil {
		genreMap := map[string]int32{}
		for _, g := range genreResp.Genres {
			genreMap[strings.ToLower(g.Name)] = g.Id
		}
		for _, gName := range req.Genres {
			if id, ok := genreMap[strings.ToLower(gName)]; ok {
				genreIDs = append(genreIDs, id)
			}
		}
	}

	novel, err := h.Clients.Novel.CreateNovel(r.Context(), &novelv1.CreateNovelRequest{
		Title:    sanitize(req.Title),
		Synopsis: sanitize(req.Synopsis),
		Author:   sanitize(req.Author),
		CoverUrl: coverURL,
		GenreIds: genreIDs,
	})
	if err != nil {
		writeError(w, 409, err.Error())
		return
	}
	writeJSON(w, 201, novel)
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

	baseURL := h.r2BaseURL()
	coverURL := req.CoverURL
	if baseURL != "" && strings.HasPrefix(coverURL, baseURL+"/") {
		coverURL = strings.TrimPrefix(coverURL, baseURL+"/")
	}

	var genreIDs []int32
	genreResp, _ := h.Clients.Novel.ListGenres(r.Context(), &novelv1.ListGenresRequest{})
	if genreResp != nil {
		genreMap := map[string]int32{}
		for _, g := range genreResp.Genres {
			genreMap[strings.ToLower(g.Name)] = g.Id
		}
		for _, gName := range req.Genres {
			if gid, ok := genreMap[strings.ToLower(gName)]; ok {
				genreIDs = append(genreIDs, gid)
			}
		}
	}

	_, err := h.Clients.Novel.UpdateNovel(r.Context(), &novelv1.UpdateNovelRequest{
		Id:       id,
		Title:    sanitize(req.Title),
		Synopsis: sanitize(req.Synopsis),
		Author:   sanitize(req.Author),
		Status:   req.Status,
		CoverUrl: coverURL,
		GenreIds: genreIDs,
	})
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) AdminDeleteNovel(w http.ResponseWriter, r *http.Request) {
	_, err := h.Clients.Novel.DeleteNovel(r.Context(), &novelv1.DeleteNovelRequest{
		Id: chi.URLParam(r, "id"),
	})
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handlers) AdminListNovels(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Clients.Novel.ListNovels(r.Context(), &novelv1.ListNovelsRequest{Page: 1, PageSize: 100})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"novels":         resp.Novels,
		"total":          resp.Total,
		"cover_base_url": h.r2BaseURL(),
	})
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
	ch, err := h.Clients.Novel.CreateChapter(r.Context(), &novelv1.CreateChapterRequest{
		NovelId: req.NovelID,
		Number:  int32(req.Number),
		Title:   sanitize(req.Title),
		Content: req.Content,
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, map[string]interface{}{"id": ch.Id, "word_count": ch.WordCount})
}

func (h *Handlers) AdminUpdateChapter(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	readJSON(r, &req)
	_, err := h.Clients.Novel.UpdateChapter(r.Context(), &novelv1.UpdateChapterRequest{
		Id:      chi.URLParam(r, "id"),
		Title:   sanitize(req.Title),
		Content: req.Content,
	})
	if err != nil {
		writeError(w, 404, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) AdminDeleteChapter(w http.ResponseWriter, r *http.Request) {
	_, err := h.Clients.Novel.DeleteChapter(r.Context(), &novelv1.DeleteChapterRequest{
		Id: chi.URLParam(r, "id"),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

func (h *Handlers) AdminListChapters(w http.ResponseWriter, r *http.Request) {
	novelsResp, _ := h.Clients.Novel.ListNovels(r.Context(), &novelv1.ListNovelsRequest{Page: 1, PageSize: 100})
	var allChapters []interface{}
	if novelsResp != nil {
		for _, novel := range novelsResp.Novels {
			chResp, err := h.Clients.Novel.ListChapters(r.Context(), &novelv1.ListChaptersRequest{
				NovelId: novel.Id, Page: 1, PageSize: 500,
			})
			if err == nil {
				for _, ch := range chResp.Chapters {
					allChapters = append(allChapters, map[string]interface{}{
						"chapter":     ch,
						"novel_title": novel.Title,
						"novel_id":    novel.Id,
						"novel_slug":  novel.Slug,
					})
				}
			}
		}
	}
	writeJSON(w, 200, map[string]interface{}{"chapters": allChapters, "total": len(allChapters)})
}

func (h *Handlers) AdminGetChapter(w http.ResponseWriter, r *http.Request) {
	novelSlug := r.URL.Query().Get("novel_slug")
	numberStr := r.URL.Query().Get("number")
	if novelSlug == "" || numberStr == "" {
		writeError(w, 400, "novel_slug and number are required")
		return
	}
	num, err := strconv.Atoi(numberStr)
	if err != nil {
		writeError(w, 400, "invalid chapter number")
		return
	}
	ch, err := h.Clients.Novel.GetChapter(r.Context(), &novelv1.GetChapterRequest{
		NovelSlug:     novelSlug,
		ChapterNumber: int32(num),
	})
	if err != nil {
		writeError(w, 404, "chapter not found")
		return
	}
	writeJSON(w, 200, map[string]interface{}{
		"id":         ch.Id,
		"novel_id":   ch.NovelId,
		"number":     ch.Number,
		"title":      ch.Title,
		"content":    ch.Content,
		"word_count": ch.WordCount,
	})
}

func (h *Handlers) AdminListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 50
	}

	resp, err := h.Clients.User.ListUsers(r.Context(), &userv1.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{"users": resp.Users, "total": resp.Total})
}

func (h *Handlers) AdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	uid := chi.URLParam(r, "id")
	currentUser := middleware.GetUserID(r.Context())
	if uid == currentUser && req.Role != "admin" {
		writeError(w, 403, "cannot demote yourself")
		return
	}

	_, err := h.Clients.User.UpdateUserRole(r.Context(), &userv1.UpdateUserRoleRequest{
		UserId: uid,
		Role:   req.Role,
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "updated"})
}

func (h *Handlers) AdminCreateGenre(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := readJSON(r, &req); err != nil || strings.TrimSpace(req.Name) == "" {
		writeError(w, 400, "name is required")
		return
	}
	resp, err := h.Clients.Novel.CreateGenre(r.Context(), &novelv1.CreateGenreRequest{
		Name: sanitize(req.Name),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, resp)
}

func (h *Handlers) AdminListGenres(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Clients.Novel.ListGenres(r.Context(), &novelv1.ListGenresRequest{})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]interface{}{"genres": resp.Genres})
}

func (h *Handlers) AdminDeleteGenre(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, 400, "invalid id")
		return
	}
	_, err = h.Clients.Novel.DeleteGenre(r.Context(), &novelv1.DeleteGenreRequest{
		Id: int32(id),
	})
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}
