package grpclog

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type requestIDKey struct{}

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		requestID := extractRequestID(ctx)
		ctx = context.WithValue(ctx, requestIDKey{}, requestID)

		defer func() {
			if r := recover(); r != nil {
				logger.Error("PANIC recovered in gRPC handler",
					zap.String("method", info.FullMethod),
					zap.String("request_id", requestID),
					zap.Any("panic", r),
					zap.Duration("duration", time.Since(start)),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		// Execute the handler
		resp, err = handler(ctx, req)

		duration := time.Since(start)
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}

		// Log based on outcome
		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.String("request_id", requestID),
			zap.String("code", code.String()),
			zap.Duration("duration", duration),
		}

		if err != nil {
			logger.Error("gRPC call failed", append(fields, zap.Error(err))...)
		} else if duration > 2*time.Second {
			logger.Warn("gRPC call slow", fields...)
		} else {
			logger.Info("gRPC call completed", fields...)
		}

		return resp, err
	}
}

// RequestIDFromContext extracts the request ID from context.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

func extractRequestID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			return vals[0]
		}
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
