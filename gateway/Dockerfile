FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY gateway/go.mod gateway/go.sum ./

RUN go mod download

COPY gateway/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gateway cmd/gateway/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/gateway /bin/gateway

EXPOSE 8080

CMD ["/bin/gateway"]
