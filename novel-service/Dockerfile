FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY proto/ /proto/
COPY pkg/ /pkg/
COPY novel-service/go.mod novel-service/go.sum* ./
RUN go mod download
COPY novel-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server cmd/server/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/server /bin/server
EXPOSE 50052
CMD ["/bin/server"]
