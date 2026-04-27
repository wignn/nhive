use search_indexer::ElasticIndexer;
use std::env;
use std::sync::Arc;
use tonic::transport::Server;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let es_url = env::var("ELASTICSEARCH_URL").unwrap_or_else(|_| "http://localhost:9200".into());
    let nats_url = env::var("NATS_URL").unwrap_or_else(|_| "nats://localhost:4222".into());
    let port = env::var("GRPC_PORT").unwrap_or_else(|_| "50054".into());

    let indexer = Arc::new(ElasticIndexer::new(&es_url).await?);

    // Start NATS subscriber in background
    let sub_indexer = indexer.clone();
    tokio::spawn(async move {
        if let Err(e) = search_subscriber::start_subscriber(&nats_url, sub_indexer).await {
            tracing::error!("NATS subscriber error: {}", e);
        }
    });

    let addr = format!("0.0.0.0:{}", port).parse()?;
    tracing::info!("Search Service listening on {}", addr);

    Server::builder()
        .add_service(tonic_reflection::server::Builder::configure().build()?)
        .serve(addr)
        .await?;

    Ok(())
}
