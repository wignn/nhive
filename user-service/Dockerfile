FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY user-service/go.mod ./
RUN go mod download

COPY user-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server cmd/server/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/server /bin/server
EXPOSE 50051
CMD ["/bin/server"]
