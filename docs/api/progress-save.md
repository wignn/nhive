# Save Progress API (Protected)

Saves the user's reading progress for a specific novel. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`PUT /api/v1/progress/{novelId}`

**Request Body (JSON):**
```json
{
  "last_chapter_number": 45,
  "last_scroll_percent": 85.5
}
```

**Response (200 OK):** 
```json
{
  "status": "saved"
}
```
