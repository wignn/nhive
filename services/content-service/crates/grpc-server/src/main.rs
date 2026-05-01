use content_storage::postgres::PostgresStorage;
use content_storage::cache::RedisCache;
use sqlx::postgres::PgPoolOptions;
use std::env;
use std::net::SocketAddr;
use tonic::transport::Server;
use tracing_subscriber::{self, EnvFilter, fmt};

mod service;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let env = env::var("APP_ENV").unwrap_or_else(|_| "development".into());
    if env == "production" {
        fmt::Subscriber::builder()
            .json()
            .with_env_filter(EnvFilter::from_default_env().add_directive("content_service=info".parse()?))
            .with_target(true)
            .with_thread_ids(true)
            .with_file(true)
            .with_line_number(true)
            .init();
    } else {
        fmt::Subscriber::builder()
            .pretty()
            .with_env_filter(EnvFilter::from_default_env().add_directive("content_service=debug".parse()?))
            .with_target(true)
            .init();
    }

    tracing::info!(service = "content-service", "starting content-service");

    let db_url = env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://novelhive:secret@localhost:5432/novelhive_novels".into());
    let redis_url = env::var("REDIS_URL")
        .unwrap_or_else(|_| "redis://localhost:6379/2".into());
    let port = env::var("GRPC_PORT").unwrap_or_else(|_| "50053".into());

    tracing::info!(service = "content-service", "connecting to database");
    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&db_url)
        .await?;
    tracing::info!(service = "content-service", "database connected");

    let storage = PostgresStorage::new(pool);
    let cache = match RedisCache::new(&redis_url).await {
        Ok(c) => {
            tracing::info!(service = "content-service", url = %redis_url, "redis cache connected");
            Some(c)
        }
        Err(e) => {
            tracing::warn!(service = "content-service", error = %e, "redis cache unavailable (non-fatal)");
            None
        }
    };

    let addr: SocketAddr = format!("0.0.0.0:{}", port).parse()?;
    tracing::info!(service = "content-service", %addr, "content-service listening");

    Server::builder()
        .add_service(tonic_reflection::server::Builder::configure().build_v1()?)
        .serve(addr)
        .await?;

    tracing::info!(service = "content-service", "content-service stopped");

    Ok(())
}
