# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy entire go-redux workspace
COPY . .

# Build from orchestrator directory
WORKDIR /build/examples/saga-order-fulfillment/orchestrator

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o orchestrator ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/examples/saga-order-fulfillment/orchestrator/orchestrator .
COPY --from=builder /build/examples/saga-order-fulfillment/orchestrator/orchestrator.json .

EXPOSE 8080

CMD ["./orchestrator"]
