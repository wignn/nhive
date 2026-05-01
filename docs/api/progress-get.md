# Get Progress API (Protected)

Fetches the user's reading progress for a specific novel. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`GET /api/v1/progress/{novelId}`

**Response (200 OK):**
```json
{
  "progress": {
    "user_id": "user-uuid",
    "novel_id": "novel-uuid",
    "last_chapter_number": 45,
    "last_scroll_percent": 85.5,
    "updated_at": "2026-04-30T10:00:00Z"
  }
}
```
