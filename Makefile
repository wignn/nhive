.PHONY: proto dev down build migrate seed clean

# Generate protobuf code for Go services
proto:
	@echo "Generating protobuf code..."
	@for dir in user novel content search comment library; do \
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
	cd user-service && go run cmd/migrate/main.go
	cd novel-service && go run cmd/migrate/main.go
	cd comment-service && go run cmd/migrate/main.go
	cd library-service && go run cmd/migrate/main.go
	@echo "Migrations complete."

# Seed sample data
seed:
	@echo "Seeding data..."
	cd scripts && go run seed.go
	@echo "Seed complete."

# Build individual services
build-gateway:
	cd gateway && go build -o bin/gateway cmd/gateway/main.go

build-user:
	cd user-service && go build -o bin/server cmd/server/main.go

build-novel:
	cd novel-service && go build -o bin/server cmd/server/main.go

build-content:
	cd content-service && cargo build --release

build-search:
	cd search-service && cargo build --release

build-comment:
	cd comment-service && go build -o bin/server cmd/server/main.go

build-library:
	cd library-service && go build -o bin/server cmd/server/main.go

# Build all Go services
build-go: build-gateway build-user build-novel build-comment build-library

# Build all Rust services
build-rust: build-content build-search

# Build everything
build-all: build-go build-rust

# Test
test:
	cd gateway && go test ./...
	cd user-service && go test ./...
	cd novel-service && go test ./...
	cd comment-service && go test ./...
	cd library-service && go test ./...
	cd content-service && cargo test
	cd search-service && cargo test

# Tidy all Go modules (including shared packages)
tidy:
	cd pkg/logger && go mod tidy
	cd pkg/config && go mod tidy
	cd pkg/grpclog && go mod tidy
	cd pkg/grpcauth && go mod tidy
	cd gateway && go mod tidy
	cd user-service && go mod tidy
	cd novel-service && go mod tidy
	cd comment-service && go mod tidy
	cd library-service && go mod tidy
