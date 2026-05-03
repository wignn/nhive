FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY proto/ /proto/
COPY pkg/ /pkg/
COPY services/notification-service/go.mod services/notification-service/go.sum* ./

RUN go mod download

COPY services/notification-service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/server cmd/server/main.go
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/server /bin/server

EXPOSE 50057

ENTRYPOINT ["/bin/server"]
