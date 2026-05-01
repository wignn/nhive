# Search Novels API

Performs a full-text search across novels.

**Endpoint:**
`GET /api/v1/search`

**Query Parameters:**
- `q`: Search query (required)
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20)

**Response (200 OK):**
```json
{
  "results": [
    {
      "id": "novel-uuid",
      "title": "The Beginning After The End",
      "slug": "the-beginning-after-the-end",
      "author": "TurtleMe",
      "cover_url": "covers/tbate.jpg",
      "highlights": {
        "title": ["The <em>Beginning</em> After The End"],
        "synopsis": ["King Grey has unrivaled strength from the <em>beginning</em>..."]
      }
    }
  ],
  "total_hits": 15,
  "took_ms": 12,
  "cover_base_url": "https://cdn.novelhive.com"
}
```
