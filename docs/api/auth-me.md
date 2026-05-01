# Get Profile API (Protected)

Returns the authenticated user's profile data. Requires the `Authorization: Bearer <token>` header.

**Endpoint:**
`GET /api/v1/auth/me`

**Response (200 OK):**
```json
{
  "id": "user-uuid",
  "username": "johndoe",
  "email": "john@example.com",
  "role": "user",
  "avatar_url": "",
  "created_at": "2026-04-30T10:00:00Z"
}
```
