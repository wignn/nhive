# Read Chapter API

Fetches the full content of a specific chapter.

**Endpoint:**
`GET /api/v1/novels/{slug}/chapters/{number}`

**Response (200 OK):**
```json
{
  "chapter": {
    "id": "chapter-uuid",
    "novel_id": "novel-uuid",
    "number": 1,
    "title": "The Awakening",
    "content": "<p>The pain was unbearable...</p>",
    "word_count": 2500,
    "created_at": "2026-04-30T10:00:00Z"
  }
}
```
