.PHONY: proto dev down build migrate seed clean

# Generate protobuf code for Go services
proto:
	@echo "Generating protobuf code..."
	@for dir in user novel content search comment library notification; do \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/$$dir/v1/$$dir.proto; \
	done
	@echo "Proto generation complete."

# Start all services
dev:
	docker-compose up --build -d

# Stop all services
down:
	docker-compose down

# Stop and remove volumes
clean:
	docker-compose down -v

# View logs (structured JSON in production)
logs:
	docker-compose logs -f

# View specific service logs
logs-%:
	docker-compose logs -f $*

# Run database migrations for all Go services
migrate:
	@echo "Running migrations..."
	cd services/user-service && go run cmd/migrate/main.go
	cd services/novel-service && go run cmd/migrate/main.go
	cd services/comment-service && go run cmd/migrate/main.go
	cd services/library-service && go run cmd/migrate/main.go
	@echo "Migrations complete."

# Seed sample data
seed:
	@echo "Seeding data..."
	cd scripts && go run seed.go
	@echo "Seed complete."

# Build individual services
build-gateway:
	cd services/gateway && go build -o bin/gateway cmd/gateway/main.go

build-user:
	cd services/user-service && go build -o bin/server cmd/server/main.go

build-novel:
	cd services/novel-service && go build -o bin/server cmd/server/main.go

build-content:
	cd services/content-service && cargo build --release

build-search:
	cd services/search-service && cargo build --release

build-comment:
	cd services/comment-service && go build -o bin/server cmd/server/main.go

build-library:
	cd services/library-service && go build -o bin/server cmd/server/main.go

build-notification:
	cd services/notification-service && go build -o bin/server cmd/server/main.go

# Build all Go services
build-go: build-gateway build-user build-novel build-comment build-library build-notification

# Build all Rust services
build-rust: build-content build-search

# Build everything
build-all: build-go build-rust

# Test
test:
	cd services/gateway && go test ./...
	cd services/user-service && go test ./...
	cd services/novel-service && go test ./...
	cd services/comment-service && go test ./...
	cd services/library-service && go test ./...
	cd services/content-service && cargo test
	cd services/search-service && cargo test

# Tidy all Go modules (including shared packages)
tidy:
	cd pkg/logger && go mod tidy
	cd pkg/config && go mod tidy
	cd pkg/grpclog && go mod tidy
	cd pkg/grpcauth && go mod tidy
	cd services/gateway && go mod tidy
	cd services/user-service && go mod tidy
	cd services/novel-service && go mod tidy
	cd services/comment-service && go mod tidy
	cd services/library-service && go mod tidy
	cd services/notification-service && go mod tidy
