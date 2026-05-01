# List Comments API

Fetches comments for a specific chapter. (Also available as public if no auth is sent, but auth is required to like/comment).

**Endpoint:**
`GET /api/v1/chapters/{chapterId}/comments`

**Response (200 OK):**
```json
{
  "comments": [
    {
      "id": "comment-uuid",
      "chapter_id": "chapter-uuid",
      "user_id": "user-uuid",
      "username": "johndoe",
      "avatar_url": "",
      "content": "What an amazing chapter!",
      "likes_count": 42,
      "created_at": "2026-04-30T10:00:00Z"
    }
  ]
}
```
