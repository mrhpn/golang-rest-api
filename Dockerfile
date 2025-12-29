# Step 1: Build Stage
FROM golang:1.24.11-alpine3.23 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Step 2: Final Image Stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Add non-root user for security
RUN adduser -D appuser

WORKDIR /

# Copy binary, migrations and env if needed (or use env vars in k8s/docker-compose)
COPY --from=builder --chown=appuser:appuser /api /api
COPY --from=builder --chown=appuser:appuser /go/bin/goose /usr/local/bin/goose
COPY --from=builder --chown=appuser:appuser /app/migrations /migrations
COPY --chown=appuser:appuser scripts/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

# Switch to non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]