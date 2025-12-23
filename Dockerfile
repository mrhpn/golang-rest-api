# Step 1: Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

# Step 2: Final Image Stage
FROM alpine:latest

# Add non-root user for security
RUN adduser -D appuser
USER appuser

WORKDIR /

# Copy binary and env if needed (or use env vars in k8s/docker-compose)
COPY --from=builder /api /api

EXPOSE 8080

ENTRYPOINT ["/api"]