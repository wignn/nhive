# Add Bookmark API (Protected)

Creates a new bookmark. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`POST /api/v1/bookmarks`

**Request Body (JSON):**
```json
{
  "novel_id": "novel-uuid",
  "chapter_number": 45,
  "paragraph_index": 12,
  "note": "Great moment!" // Optional
}
```

**Response (200 OK):** 
```json
{
  "status": "added", 
  "bookmark_id": "bookmark-uuid"
}
```
