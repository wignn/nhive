# Google OAuth Login API

Signs in with a Google ID token. If the verified email already exists, the existing user is used. If the email is new, a reader account is created and `avatar_url` is initialized from the Google profile picture.

`POST /api/v1/auth/google`

Headers:

- `x-api-key: <gateway-api-key>`
- `Content-Type: application/json`

Request:

```json
{
  "id_token": "google-id-token"
}
```

Response:

```json
{
  "user_id": "user-id",
  "access_token": "jwt",
  "refresh_token": "jwt",
  "user": {
    "id": "user-id",
    "username": "johndoe",
    "email": "john@example.com",
    "avatar_url": "https://lh3.googleusercontent.com/...",
    "role": "reader",
    "created_at": "2026-05-16T10:00:00Z"
  }
}
```

Environment:

- Gateway: `GOOGLE_CLIENT_ID` or comma-separated `GOOGLE_CLIENT_IDS`
- Web: `PUBLIC_GOOGLE_CLIENT_ID`
- Mobile build: `--dart-define=GOOGLE_WEB_CLIENT_ID=<web-client-id>`
