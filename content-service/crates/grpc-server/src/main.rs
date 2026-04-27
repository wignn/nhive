use content_storage::postgres::PostgresStorage;
use content_storage::cache::RedisCache;
use sqlx::postgres::PgPoolOptions;
use std::env;
use std::net::SocketAddr;
use tonic::transport::Server;
use tracing_subscriber;

mod service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let db_url = env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://novelhive:secret@localhost:5432/novelhive_novels".into());
    let redis_url = env::var("REDIS_URL")
        .unwrap_or_else(|_| "redis://localhost:6379/2".into());
    let port = env::var("GRPC_PORT").unwrap_or_else(|_| "50053".into());

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&db_url)
        .await?;

    let storage = PostgresStorage::new(pool);
    let cache = RedisCache::new(&redis_url).await.ok();

    let addr: SocketAddr = format!("0.0.0.0:{}", port).parse()?;
    tracing::info!("Content Service listening on {}", addr);

    // TODO: Register tonic service once proto is compiled
    // For now, start an empty server
    // Register tonic-reflection service so clients can query reflection metadata.
    Server::builder()
        .add_service(tonic_reflection::server::Builder::configure().build()?)
        .serve(addr)
        .await?;

    Ok(())
}
