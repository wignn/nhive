FROM lukemathwalker/cargo-chef:latest-rust-1.94-bookworm AS chef
WORKDIR /app

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    protobuf-compiler \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

FROM chef AS planner

COPY services/search-service/Cargo.toml Cargo.toml
COPY services/search-service/crates crates

RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder

COPY --from=planner /app/recipe.json recipe.json

RUN cargo chef cook --release --recipe-path recipe.json

COPY services/search-service/ .

RUN cargo build --release --package search-grpc-server

RUN strip target/release/search-service || true

FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/* \
 && useradd -r -u 10001 appuser

WORKDIR /app

COPY --from=builder /app/target/release/search-service /usr/local/bin/search-service

USER appuser

EXPOSE 50054

ENTRYPOINT ["/usr/local/bin/search-service"]