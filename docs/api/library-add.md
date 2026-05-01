# Add to Library API (Protected)

Adds a novel to the user's library. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`POST /api/v1/library/{novelId}`

**Request Body (JSON):**
```json
{
  "status": "reading" // Optional: 'reading', 'plan_to_read', 'completed', 'dropped'
}
```

**Response (200 OK):** 
```json
{
  "status": "added"
}
```
