# Autocomplete Search API

Returns quick suggestions for the search bar.

**Endpoint:**
`GET /api/v1/search/autocomplete`

**Query Parameters:**
- `q`: Search query prefix (required)

**Response (200 OK):**
```json
{
  "suggestions": [
    {
      "id": "novel-uuid",
      "title": "The Beginning After The End",
      "slug": "the-beginning-after-the-end"
    }
  ],
  "took_ms": 5
}
```
