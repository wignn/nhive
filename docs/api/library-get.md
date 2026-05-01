# Get Library API (Protected)

Fetches the user's library. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`GET /api/v1/library`

**Response (200 OK):**
```json
{
  "items": [
    {
      "id": "library-item-uuid",
      "user_id": "user-uuid",
      "novel_id": "novel-uuid",
      "status": "reading",
      "added_at": "2026-04-30T10:00:00Z"
    }
  ]
}
```
