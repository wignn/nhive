use content_domain::ChapterContent;
use redis::AsyncCommands;
use std::sync::Arc;
use tokio::sync::Mutex;

pub struct RedisCache {
    client: Arc<Mutex<redis::aio::MultiplexedConnection>>,
    ttl_secs: u64,
}

impl RedisCache {
    pub async fn new(redis_url: &str) -> anyhow::Result<Self> {
        let client = redis::Client::open(redis_url)?;
        let conn = client.get_multiplexed_async_connection().await?;
        Ok(Self {
            client: Arc::new(Mutex::new(conn)),
            ttl_secs: 900, // 15 min
        })
    }

    pub async fn get_chapter(&self, novel_slug: &str, number: i32) -> Option<ChapterContent> {
        let key = format!("content:{}:{}", novel_slug, number);
        let mut conn = self.client.lock().await;
        let data: Option<String> = conn.get(&key).await.ok()?;
        data.and_then(|d| serde_json::from_str(&d).ok())
    }

    pub async fn set_chapter(&self, novel_slug: &str, chapter: &ChapterContent) {
        let key = format!("content:{}:{}", novel_slug, chapter.number);
        if let Ok(data) = serde_json::to_string(chapter) {
            let mut conn = self.client.lock().await;
            let _: Result<(), _> = conn.set_ex(&key, &data, self.ttl_secs).await;
        }
    }
}
