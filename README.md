# NovelHive

A microservice-based novel reading platform built with Go, Rust, gRPC, and modern infrastructure.

## Architecture

- **API Gateway** (Go) — REST endpoints, JWT auth, rate limiting
- **User Service** (Go) — Authentication & user management
- **Novel Service** (Go) — Novel & chapter CRUD
- **Content Service** (Rust) — High-performance content delivery
- **Search Service** (Rust) — Elasticsearch-powered full-text search
- **Comment Service** (Go) — Threaded comments with likes
- **Library Service** (Go) — Bookmarks, reading progress, reading lists
- **Frontend** (Next.js) — Immersive reading experience

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Languages | Go 1.24+, Rust 1.85+ |
| Frontend | Next.js 15, Vanilla CSS |
| Communication | gRPC (internal), REST (external) |
| Event Bus | NATS JetStream |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Search | Elasticsearch 8.13 |
| Logging | Zap (Go), tracing-subscriber (Rust) |
| Containers | Docker, Docker Compose |

## Quick Start

```bash
# Start all infrastructure and services
make dev

# View logs (structured JSON in production)
make logs

# Stop everything
make down

# Stop and clean volumes
make clean

# Tidy all Go modules
make tidy
```

## Observability

All services emit **structured JSON logs** in production mode:

```json
{
  "level": "info",
  "ts": "2026-04-29T12:00:00Z",
  "caller": "server/main.go:45",
  "msg": "user-service started",
  "service": "user-service",
  "port": "50051"
}
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | `production` for JSON logs, `development` for console | `development` |
| `LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` | `info` |
| `RUST_LOG` | Rust log filter (e.g., `search_service=info`) | — |

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Login |
| GET | `/novels` | List novels |
| GET | `/novels/:slug` | Novel detail |
| GET | `/novels/:slug/chapters/:num` | Read chapter |
| GET | `/search?q=` | Search novels |
| GET | `/chapters/:id/comments` | List comments |
| POST | `/chapters/:id/comments` | Post comment |
| GET | `/library` | User's reading list |
| PUT | `/progress/:novelId` | Save reading progress |

## Project Structure

```
novelhive/
├── pkg/                        # Shared Go packages
│   ├── logger/                 #   Structured logging (zap)
│   ├── config/                 #   Config helpers (GetEnv, MustEnv)
│   └── grpclog/                #   gRPC logging interceptor
├── proto/                      # Shared Protobuf definitions
├── gateway/                    # [Go] API Gateway
│   ├── cmd/gateway/main.go
│   ├── internal/
│   │   ├── clients/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── storage/
│   │   └── store/
│   └── Dockerfile
├── user-service/               # [Go] Auth & Users
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   ├── grpc/
│   │   ├── repository/
│   │   └── usecase/
│   └── Dockerfile
├── novel-service/              # [Go] Novels & Chapters
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/
│   │   ├── events/
│   │   ├── grpc/
│   │   ├── repository/
│   │   └── usecase/
│   └── Dockerfile
├── comment-service/            # [Go] Comments
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/
│   │   ├── grpc/
│   │   └── repository/
│   └── Dockerfile
├── library-service/            # [Go] Library & Progress
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/
│   │   ├── grpc/
│   │   └── repository/
│   └── Dockerfile
├── content-service/            # [Rust] Content Delivery
│   ├── crates/
│   │   ├── domain/
│   │   ├── grpc-server/
│   │   └── storage/
│   ├── Cargo.toml
│   └── Dockerfile
├── search-service/             # [Rust] Search Engine
│   ├── crates/
│   │   ├── domain/
│   │   ├── grpc-server/
│   │   ├── indexer/
│   │   └── subscriber/
│   ├── Cargo.toml
│   └── Dockerfile
├── frontend/                   # [Next.js] Web UI
├── scripts/                    # DB init, seed data
├── docker-compose.yml
└── Makefile
```
