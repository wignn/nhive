module github.com/novelhive/novel-service

go 1.24

require (
	github.com/gosimple/slug v1.14.0
	github.com/jackc/pgx/v5 v5.6.0
	github.com/nats-io/nats.go v1.35.0
	github.com/novelhive/pkg/grpcauth v0.0.0
	github.com/novelhive/pkg/grpclog v0.0.0
	github.com/novelhive/pkg/logger v0.0.0
	github.com/novelhive/proto v0.0.0
	github.com/redis/go-redis/v9 v9.5.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.71.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250219182151-9fdb1cabc7b2 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace (
	github.com/novelhive/pkg/grpcauth => ../pkg/grpcauth
	github.com/novelhive/pkg/grpclog => ../pkg/grpclog
	github.com/novelhive/pkg/logger => ../pkg/logger
	github.com/novelhive/proto => ../proto
)
