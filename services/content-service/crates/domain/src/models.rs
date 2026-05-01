use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChapterContent {
    pub id: String,
    pub novel_id: String,
    pub novel_title: String,
    pub novel_slug: String,
    pub number: i32,
    pub title: String,
    pub content: String,
    pub word_count: i32,
    pub total_chapters: i32,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AdjacentChapters {
    pub has_prev: bool,
    pub has_next: bool,
    pub prev_number: i32,
    pub next_number: i32,
}
