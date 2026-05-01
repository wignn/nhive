# Create Comment API (Protected)

Posts a new comment on a chapter. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`POST /api/v1/chapters/{chapterId}/comments`

**Request Body (JSON):**
```json
{
  "content": "What an amazing chapter!"
}
```

**Response (200 OK):**
```json
{
  "comment": {
    "id": "comment-uuid",
    "chapter_id": "chapter-uuid",
    // ... comment details
  }
}
```
