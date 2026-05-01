# Remove from Library API (Protected)

Removes a novel from the user's library. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`DELETE /api/v1/library/{novelId}`

**Response (200 OK):** 
```json
{
  "status": "removed"
}
```
