# Register API

Creates a new user account.

**Endpoint:**
`POST /api/v1/auth/register`

**Request Body (JSON):**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```
*Validation:* Username ≥ 3 chars, Password ≥ 6 chars.

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
