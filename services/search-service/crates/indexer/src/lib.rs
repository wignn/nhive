use elasticsearch::{Elasticsearch, http::transport::Transport, IndexParts, SearchParts, DeleteParts};
use search_domain::*;
use serde_json::{json, Value};

pub struct ElasticIndexer {
    client: Elasticsearch,
    index_name: String,
}

impl ElasticIndexer {
    pub async fn new(es_url: &str) -> anyhow::Result<Self> {
        let transport = Transport::single_node(es_url)?;
        let client = Elasticsearch::new(transport);
        let indexer = Self { client, index_name: "novels".to_string() };
        indexer.ensure_index().await?;
        Ok(indexer)
    }

    async fn ensure_index(&self) -> anyhow::Result<()> {
        let exists = self.client
            .indices()
            .exists(elasticsearch::indices::IndicesExistsParts::Index(&[&self.index_name]))
            .send().await?;

        if exists.status_code().as_u16() == 404 {
            let mapping = json!({
                "mappings": {
                    "properties": {
                        "title": { "type": "text", "analyzer": "standard", "fields": { "suggest": { "type": "completion" } } },
                        "synopsis": { "type": "text" },
                        "author": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
                        "genres": { "type": "keyword" },
                        "status": { "type": "keyword" },
                        "slug": { "type": "keyword" },
                        "cover_url": { "type": "keyword", "index": false },
                        "total_chapters": { "type": "integer" }
                    }
                }
            });
            self.client.indices()
                .create(elasticsearch::indices::IndicesCreateParts::Index(&self.index_name))
                .body(mapping).send().await?;
            tracing::info!("Created index: {}", self.index_name);
        }
        Ok(())
    }

    pub async fn index_novel(&self, doc: &NovelDocument) -> anyhow::Result<()> {
        self.client.index(IndexParts::IndexId(&self.index_name, &doc.id))
            .body(json!({
                "title": doc.title, "slug": doc.slug, "synopsis": doc.synopsis,
                "author": doc.author, "cover_url": doc.cover_url, "genres": doc.genres,
                "status": doc.status, "total_chapters": doc.total_chapters
            }))
            .send().await?;
        Ok(())
    }

    pub async fn remove_novel(&self, id: &str) -> anyhow::Result<()> {
        self.client.delete(DeleteParts::IndexId(&self.index_name, id)).send().await?;
        Ok(())
    }

    pub async fn search(&self, params: &SearchParams) -> anyhow::Result<SearchResult> {
        let from = ((params.page - 1).max(0)) * params.page_size;
        let mut must = vec![json!({
            "multi_match": {
                "query": params.query,
                "fields": ["title^3", "synopsis", "author^2"],
                "fuzziness": "AUTO"
            }
        })];

        if let Some(ref genre) = params.genre {
            must.push(json!({ "term": { "genres": genre } }));
        }
        if let Some(ref status) = params.status {
            must.push(json!({ "term": { "status": status } }));
        }

        let body = json!({
            "from": from, "size": params.page_size,
            "query": { "bool": { "must": must } },
            "highlight": { "fields": { "title": {}, "synopsis": {} } }
        });

        let response = self.client.search(SearchParts::Index(&[&self.index_name]))
            .body(body).send().await?;
        let json: Value = response.json().await?;

        let took = json["took"].as_f64().unwrap_or(0.0) as f32;
        let total = json["hits"]["total"]["value"].as_i64().unwrap_or(0) as i32;
        let mut hits = Vec::new();

        if let Some(arr) = json["hits"]["hits"].as_array() {
            for hit in arr {
                let src = &hit["_source"];
                let highlight_text = hit["highlight"]["title"].as_array()
                    .and_then(|a| a.first())
                    .and_then(|v| v.as_str())
                    .unwrap_or("")
                    .to_string();

                hits.push(SearchHit {
                    id: hit["_id"].as_str().unwrap_or("").to_string(),
                    title: src["title"].as_str().unwrap_or("").to_string(),
                    slug: src["slug"].as_str().unwrap_or("").to_string(),
                    synopsis: src["synopsis"].as_str().unwrap_or("").to_string(),
                    author: src["author"].as_str().unwrap_or("").to_string(),
                    cover_url: src["cover_url"].as_str().unwrap_or("").to_string(),
                    genres: src["genres"].as_array()
                        .map(|a| a.iter().filter_map(|v| v.as_str().map(String::from)).collect())
                        .unwrap_or_default(),
                    score: hit["_score"].as_f64().unwrap_or(0.0) as f32,
                    highlight: highlight_text,
                });
            }
        }

        Ok(SearchResult { hits, total, page: params.page, took_ms: took })
    }

    pub async fn autocomplete(&self, query: &str, limit: i32) -> anyhow::Result<Vec<AutocompleteSuggestion>> {
        let body = json!({
            "suggest": {
                "novel-suggest": {
                    "prefix": query,
                    "completion": { "field": "title.suggest", "size": limit, "fuzzy": { "fuzziness": 1 } }
                }
            }
        });
        let response = self.client.search(SearchParts::Index(&[&self.index_name]))
            .body(body).send().await?;
        let json: Value = response.json().await?;

        let mut suggestions = Vec::new();
        if let Some(opts) = json["suggest"]["novel-suggest"].as_array()
            .and_then(|a| a.first())
            .and_then(|v| v["options"].as_array())
        {
            for opt in opts {
                let src = &opt["_source"];
                suggestions.push(AutocompleteSuggestion {
                    title: src["title"].as_str().unwrap_or("").to_string(),
                    slug: src["slug"].as_str().unwrap_or("").to_string(),
                    author: src["author"].as_str().unwrap_or("").to_string(),
                });
            }
        }
        Ok(suggestions)
    }
}
