# Multi-stage build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy entire monorepo
COPY . .

# Build from event-sourced-wallet directory
WORKDIR /build/examples/event-sourced-wallet

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/examples/event-sourced-wallet/wallet .

EXPOSE 8080

CMD ["./wallet"]
