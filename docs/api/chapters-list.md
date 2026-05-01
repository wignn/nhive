# List Chapters API

Fetches a paginated list of chapter summaries (without content) for a specific novel.

**Endpoint:**
`GET /api/v1/novels/{slug}/chapters`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 100, max: 500)

**Response (200 OK):**
```json
{
  "chapters": [
    {
      "id": "chapter-uuid",
      "number": 1,
      "title": "The Awakening",
      "word_count": 2500,
      "created_at": "2026-04-30T10:00:00Z"
    }
  ],
  "total": 450,
  "page": 1,
  "page_size": 100
}
```
