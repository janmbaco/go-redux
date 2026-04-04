# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy root go.mod
COPY go.mod go.sum ./

# Copy go-redux library source
COPY actions/ ./actions/
COPY events/ ./events/
COPY handlers/ ./handlers/
COPY ioc/ ./ioc/
COPY reducers/ ./reducers/
COPY store.go ./

# Copy example go.mod
COPY examples/workflow-employee-onboarding/go.mod examples/workflow-employee-onboarding/go.sum ./examples/workflow-employee-onboarding/

# Copy example source
COPY examples/workflow-employee-onboarding/ ./examples/workflow-employee-onboarding/

# Build from example directory
WORKDIR /build/examples/workflow-employee-onboarding
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /onboarding ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /onboarding .

EXPOSE 8080

CMD ["./onboarding"]
