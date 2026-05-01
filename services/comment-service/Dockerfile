FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY proto/ /proto/

COPY pkg/ /pkg/

COPY services/comment-service/go.mod services/comment-service/go.sum* ./

RUN go mod download

COPY services/comment-service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /bin/server cmd/server/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/server /bin/server

EXPOSE 50055

ENTRYPOINT ["/bin/server"]