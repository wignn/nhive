# List Bookmarks API (Protected)

Fetches the user's bookmarks. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`GET /api/v1/bookmarks`

**Response (200 OK):**
```json
{
  "bookmarks": [
    {
      "id": "bookmark-uuid",
      "user_id": "user-uuid",
      "novel_id": "novel-uuid",
      "chapter_number": 45,
      "paragraph_index": 12,
      "note": "Great moment!",
      "created_at": "2026-04-30T10:00:00Z"
    }
  ]
}
```
