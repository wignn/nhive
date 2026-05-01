use async_nats::jetstream;
use futures_util::TryStreamExt;
use search_domain::NovelDocument;
use search_indexer::ElasticIndexer;
use std::sync::Arc;

pub async fn start_subscriber(
    nats_url: &str,
    indexer: Arc<ElasticIndexer>,
) -> anyhow::Result<()> {
    let client = async_nats::connect(nats_url).await?;
    let js = jetstream::new(client);

    // Get or create stream
    let stream = js
        .get_or_create_stream(jetstream::stream::Config {
            name: "NOVELS".to_string(),
            subjects: vec![
                "novel.>".to_string(),
                "chapter.>".to_string(),
            ],
            ..Default::default()
        })
        .await?;

    // Durable pull consumer
    let consumer = stream
        .get_or_create_consumer(
            "search-indexer",
            jetstream::consumer::pull::Config {
                durable_name: Some("search-indexer".to_string()),
                filter_subject: "novel.>".to_string(),
                ..Default::default()
            },
        )
        .await?;

    tracing::info!("NATS subscriber started, listening for novel events...");

    loop {
        let mut messages = consumer
            .fetch()
            .max_messages(10)
            .messages()
            .await?;

        while let Some(msg) = messages.try_next().await.map_err(anyhow::Error::msg)? {
            let subject = msg.subject.as_str();

            match subject {
                "novel.created" | "novel.updated" => {
                    match serde_json::from_slice::<NovelDocument>(&msg.payload) {
                        Ok(doc) => {
                            match indexer.index_novel(&doc).await {
                                Ok(_) => {
                                    tracing::info!(
                                        "Indexed novel: {} ({})",
                                        doc.title,
                                        subject
                                    );
                                }
                                Err(err) => {
                                    tracing::error!(
                                        "Failed to index novel {}: {}",
                                        doc.id,
                                        err
                                    );
                                }
                            }
                        }
                        Err(err) => {
                            tracing::error!(
                                "Invalid payload for {}: {}",
                                subject,
                                err
                            );
                        }
                    }
                }

                _ => {
                    tracing::debug!("Ignored event: {}", subject);
                }
            }

            if let Err(err) = msg.ack().await {
                tracing::error!("Failed to ack message: {}", err);
            }
        }

        tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
    }
}