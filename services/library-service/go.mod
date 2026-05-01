module github.com/novelhive/library-service

go 1.24

require (
	github.com/jackc/pgx/v5 v5.6.0
	github.com/novelhive/pkg/grpcauth v0.0.0
	github.com/novelhive/pkg/grpclog v0.0.0
	github.com/novelhive/pkg/logger v0.0.0
	github.com/novelhive/proto v0.0.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.71.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
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
	github.com/novelhive/pkg/grpcauth => ../../pkg/grpcauth
	github.com/novelhive/pkg/grpclog => ../../pkg/grpclog
	github.com/novelhive/pkg/logger => ../../pkg/logger
	github.com/novelhive/proto => ../../proto
)
