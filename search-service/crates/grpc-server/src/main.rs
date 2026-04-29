use search_indexer::ElasticIndexer;
use std::env;
use std::sync::Arc;
use tonic::transport::Server;
use tracing_subscriber::{self, EnvFilter, fmt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Structured logging: JSON in production, pretty console in development
    let env_mode = env::var("APP_ENV").unwrap_or_else(|_| "development".into());
    if env_mode == "production" {
        fmt::Subscriber::builder()
            .json()
            .with_env_filter(EnvFilter::from_default_env().add_directive("search_service=info".parse()?))
            .with_target(true)
            .with_thread_ids(true)
            .with_file(true)
            .with_line_number(true)
            .init();
    } else {
        fmt::Subscriber::builder()
            .pretty()
            .with_env_filter(EnvFilter::from_default_env().add_directive("search_service=debug".parse()?))
            .with_target(true)
            .init();
    }

    tracing::info!(service = "search-service", "starting search-service");

    let es_url = env::var("ELASTICSEARCH_URL").unwrap_or_else(|_| "http://localhost:9200".into());
    let nats_url = env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".into());
    let port = env::var("GRPC_PORT").unwrap_or_else(|_| "50054".into());

    tracing::info!(service = "search-service", url = %es_url, "connecting to elasticsearch");
    let indexer = Arc::new(ElasticIndexer::new(&es_url).await?);
    tracing::info!(service = "search-service", "elasticsearch connected");

    // Start NATS subscriber in background
    let sub_indexer = indexer.clone();
    let nats_url_clone = nats_url.clone();
    tokio::spawn(async move {
        tracing::info!(service = "search-service", url = %nats_url_clone, "starting NATS subscriber");
        if let Err(e) = search_subscriber::start_subscriber(&nats_url_clone, sub_indexer).await {
            tracing::error!(service = "search-service", error = %e, "NATS subscriber failed");
        }
    });

    let addr = format!("0.0.0.0:{}", port).parse()?;
    tracing::info!(service = "search-service", %addr, "search-service listening");

    Server::builder()
        .add_service(tonic_reflection::server::Builder::configure().build()?)
        .serve(addr)
        .await?;

    tracing::info!(service = "search-service", "search-service stopped");

    Ok(())
}
