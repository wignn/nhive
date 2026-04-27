# NovelHive 📚

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
| Languages | Go 1.22+, Rust 1.77+ |
| Frontend | Next.js 15, Vanilla CSS |
| Communication | gRPC (internal), REST (external) |
| Event Bus | NATS JetStream |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Search | Elasticsearch 8.13 |
| Containers | Docker, Docker Compose |

## Quick Start

```bash
# Start all infrastructure and services
make dev

# View logs
make logs

# Stop everything
make down

# Stop and clean volumes
make clean
```

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
n/
├── proto/              # Shared Protobuf definitions
├── gateway/            # [Go] API Gateway
├── user-service/       # [Go] Auth & Users
├── novel-service/      # [Go] Novels & Chapters
├── content-service/    # [Rust] Content Delivery
├── search-service/     # [Rust] Search Engine
├── comment-service/    # [Go] Comments
├── library-service/    # [Go] Library & Progress
├── frontend/           # [Next.js] Web UI
├── scripts/            # DB init, seed data
├── docker-compose.yml
└── Makefile
```
