# Update Library Status API (Protected)

Updates the reading status of a novel in the library. Requires `Authorization: Bearer <token>`.

**Endpoint:**
`PUT /api/v1/library/{novelId}/status`

**Request Body (JSON):**
```json
{
  "status": "completed"
}
```

**Response (200 OK):** 
```json
{
  "status": "updated"
}
```
