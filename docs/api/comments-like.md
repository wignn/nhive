# Like Comment API (Protected)

Likes or unlikes a comment (toggle). Requires `Authorization: Bearer <token>`.

**Endpoint:**
`POST /api/v1/comments/{commentId}/like`

**Response (200 OK):**
```json
{
  "status": "liked" 
  // or {"status": "unliked"}
}
```
