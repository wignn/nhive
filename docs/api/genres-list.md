# List Genres API

Fetches all available genres.

**Endpoint:**
`GET /api/v1/genres`

**Response (200 OK):**
```json
{
  "genres": [
    {"id": 1, "name": "Fantasy", "slug": "fantasy"},
    {"id": 2, "name": "Action", "slug": "action"}
  ]
}
```
