# Upload Profile Avatar API (Protected)

Uploads the authenticated user's profile photo to Cloudflare R2 and stores the resulting public URL on the user profile.

`POST /api/v1/auth/avatar`

Headers:

- `Authorization: Bearer <token>`
- `x-api-key: <gateway-api-key>`
- `Content-Type: multipart/form-data`

Form fields:

- `image`: JPEG, PNG, WebP, or GIF image, max 5MB. `avatar` is also accepted as a fallback field name.

Response:

```json
{
  "avatar_url": "https://cdn.example.com/avatars/user-id/123456789.png",
  "user": {
    "id": "user-id",
    "username": "johndoe",
    "email": "john@example.com",
    "avatar_url": "https://cdn.example.com/avatars/user-id/123456789.png",
    "role": "reader",
    "created_at": "2026-05-16T10:00:00Z"
  }
}
```
