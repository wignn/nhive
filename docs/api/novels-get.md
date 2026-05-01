# Get Novel Details API

Fetches full details for a specific novel.

**Endpoint:**
`GET /api/v1/novels/{slug}`

**Response (200 OK):**
```json
{
  "novel": {
    "id": "novel-uuid",
    "title": "The Beginning After The End",
    // ... other novel fields
  },
  "cover_base_url": "https://cdn.novelhive.com"
}
```
