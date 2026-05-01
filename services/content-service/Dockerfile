FROM lukemathwalker/cargo-chef:latest-rust-1.94-bookworm AS chef

WORKDIR /app

RUN apt-get update && apt-get install -y \
    pkg-config \
    protobuf-compiler \
    libssl-dev \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

 FROM chef AS planner

COPY services/content-service/Cargo.toml Cargo.toml
COPY services/content-service/crates crates

RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder

COPY --from=planner /app/recipe.json recipe.json

RUN cargo chef cook --release --recipe-path recipe.json

COPY services/content-service/ .

RUN cargo build --release --package content-grpc-server

RUN strip target/release/content-service || true

FROM debian:bookworm-slim AS runtime

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/* \
 && useradd -r -u 10001 appuser

WORKDIR /app

COPY --from=builder /app/target/release/content-service /usr/local/bin/content-service

USER appuser

EXPOSE 50053

ENTRYPOINT ["/usr/local/bin/content-service"]