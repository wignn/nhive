# Login API

Authenticates an existing user and returns a JWT.

**Endpoint:**
`POST /api/v1/auth/login`

**Request Body (JSON):**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5c...",
  "user": {
    "id": "user-uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user",
    "avatar_url": ""
  }
}
```
