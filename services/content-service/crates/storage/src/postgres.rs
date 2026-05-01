use content_domain::{ChapterContent, AdjacentChapters, ContentError};
use sqlx::PgPool;

pub struct PostgresStorage {
    pool: PgPool,
}

impl PostgresStorage {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    pub async fn get_chapter_content(&self, novel_slug: &str, chapter_number: i32) -> Result<ChapterContent, ContentError> {
        let row = sqlx::query_as::<_, ChapterRow>(
            r#"
            SELECT c.id, c.novel_id, n.title as novel_title, n.slug as novel_slug,
                   c.number, c.title, c.content, c.word_count, n.total_chapters, c.created_at
            FROM chapters c
            JOIN novels n ON c.novel_id = n.id
            WHERE n.slug = $1 AND c.number = $2
            "#,
        )
        .bind(novel_slug)
        .bind(chapter_number)
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| ContentError::Database(e.to_string()))?;

        match row {
            Some(r) => Ok(ChapterContent {
                id: r.id,
                novel_id: r.novel_id,
                novel_title: r.novel_title,
                novel_slug: r.novel_slug,
                number: r.number,
                title: r.title,
                content: r.content,
                word_count: r.word_count,
                total_chapters: r.total_chapters,
                created_at: r.created_at,
            }),
            None => Err(ContentError::ChapterNotFound),
        }
    }

    pub async fn get_adjacent_chapters(&self, novel_id: &str, current: i32) -> Result<AdjacentChapters, ContentError> {
        let prev_exists = sqlx::query_scalar::<_, bool>(
            "SELECT EXISTS(SELECT 1 FROM chapters WHERE novel_id = $1 AND number = $2)"
        )
        .bind(novel_id)
        .bind(current - 1)
        .fetch_one(&self.pool)
        .await
        .unwrap_or(false);

        let next_exists = sqlx::query_scalar::<_, bool>(
            "SELECT EXISTS(SELECT 1 FROM chapters WHERE novel_id = $1 AND number = $2)"
        )
        .bind(novel_id)
        .bind(current + 1)
        .fetch_one(&self.pool)
        .await
        .unwrap_or(false);

        Ok(AdjacentChapters {
            has_prev: prev_exists,
            has_next: next_exists,
            prev_number: current - 1,
            next_number: current + 1,
        })
    }
}

#[derive(sqlx::FromRow)]
struct ChapterRow {
    id: String,
    novel_id: String,
    novel_title: String,
    novel_slug: String,
    number: i32,
    title: String,
    content: String,
    word_count: i32,
    total_chapters: i32,
    created_at: chrono::DateTime<chrono::Utc>,
}
