# List Novels API

Fetches a paginated list of novels.

**Endpoint:**
`GET /api/v1/novels`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20, max: 50)
- `genre` (optional): Filter by genre slug
- `status` (optional): Filter by status (`ongoing`, `completed`, `hiatus`)
- `sort` (optional): Sort order (`newest`, `popular`, `alphabetical`)

**Response (200 OK):**
```json
{
  "novels": [
    {
      "id": "novel-uuid",
      "title": "The Beginning After The End",
      "slug": "the-beginning-after-the-end",
      "synopsis": "King Grey has unrivaled strength...",
      "cover_url": "covers/tbate.jpg",
      "author": "TurtleMe",
      "status": "ongoing",
      "total_chapters": 450,
      "genres": [{"id": 1, "name": "Fantasy", "slug": "fantasy"}],
      "created_at": "...",
      "updated_at": "..."
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "cover_base_url": "https://cdn.novelhive.com"
}
```
