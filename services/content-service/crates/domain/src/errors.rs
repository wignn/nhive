use thiserror::Error;

#[derive(Debug, Error)]
pub enum ContentError {
    #[error("Chapter not found")]
    ChapterNotFound,
    #[error("Novel not found")]
    NovelNotFound,
    #[error("Database error: {0}")]
    Database(String),
    #[error("Cache error: {0}")]
    Cache(String),
}
