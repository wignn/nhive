FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY proto/ /proto/

COPY pkg/ /pkg/

COPY services/gateway/go.mod services/gateway/go.sum ./

RUN go mod download

COPY services/gateway/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/gateway cmd/gateway/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/gateway /bin/gateway

EXPOSE 8080

ENTRYPOINT ["/bin/gateway"]