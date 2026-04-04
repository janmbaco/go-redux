# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy entire go-redux workspace
COPY . .

# Build from payment-service directory
WORKDIR /build/examples/saga-order-fulfillment/services/payment-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/examples/saga-order-fulfillment/services/payment-service/payment-service .
COPY --from=builder /build/examples/saga-order-fulfillment/services/payment-service/payment-service.json .

EXPOSE 8081

CMD ["./payment-service"]
