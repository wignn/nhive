use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NovelDocument {
    pub id: String,
    pub title: String,
    pub slug: String,
    pub synopsis: String,
    pub author: String,
    pub cover_url: String,
    pub genres: Vec<String>,
    pub status: String,
    pub total_chapters: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SearchHit {
    pub id: String,
    pub title: String,
    pub slug: String,
    pub synopsis: String,
    pub author: String,
    pub cover_url: String,
    pub genres: Vec<String>,
    pub score: f32,
    pub highlight: String,
}

#[derive(Debug, Clone)]
pub struct SearchParams {
    pub query: String,
    pub page: i32,
    pub page_size: i32,
    pub genre: Option<String>,
    pub status: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SearchResult {
    pub hits: Vec<SearchHit>,
    pub total: i32,
    pub page: i32,
    pub took_ms: f32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AutocompleteSuggestion {
    pub title: String,
    pub slug: String,
    pub author: String,
}

#[derive(Debug, thiserror::Error)]
pub enum SearchError {
    #[error("Elasticsearch error: {0}")]
    Elasticsearch(String),
    #[error("Index not found")]
    IndexNotFound,
}
