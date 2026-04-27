use content_storage::postgres::PostgresStorage;
use content_storage::cache::RedisCache;
use std::sync::Arc;

pub struct ContentServiceImpl {
    storage: Arc<PostgresStorage>,
    cache: Option<Arc<RedisCache>>,
}

impl ContentServiceImpl {
    pub fn new(storage: PostgresStorage, cache: Option<RedisCache>) -> Self {
        Self {
            storage: Arc::new(storage),
            cache: cache.map(Arc::new),
        }
    }

    pub async fn get_chapter_content(&self, novel_slug: &str, chapter_number: i32) -> Result<content_domain::ChapterContent, content_domain::ContentError> {
        // Check cache first
        if let Some(ref cache) = self.cache {
            if let Some(cached) = cache.get_chapter(novel_slug, chapter_number).await {
                return Ok(cached);
            }
        }

        let content = self.storage.get_chapter_content(novel_slug, chapter_number).await?;

        // Cache result
        if let Some(ref cache) = self.cache {
            cache.set_chapter(novel_slug, &content).await;
        }

        Ok(content)
    }

    pub async fn get_adjacent_chapters(&self, novel_id: &str, current: i32) -> Result<content_domain::AdjacentChapters, content_domain::ContentError> {
        self.storage.get_adjacent_chapters(novel_id, current).await
    }
}
