# Docker Setup Guide

This guide explains how to run the application using Docker Compose with all
services (PostgreSQL, MinIO, Redis) configured and ready.

## Quick Start

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## Services

### 1. API Service

- **Port**: 8080
- **Health**: http://localhost:8080/health
- **Swagger**: http://localhost:8080/swagger/index.html

### 2. PostgreSQL Database

- **Port**: 5432
- **User**: `postgres` (or from `DB_USER` env)
- **Database**: `go_rest_api_db` (or from `DB_NAME` env)
- **Persistent**: Data stored in `db_data` volume

### 3. MinIO Storage

- **API Port**: 9000
- **Console Port**: 9001
- **Console**: http://localhost:9001
- **Default Credentials**: `minioadmin` / `minioadmin`
- **Persistent**: Data stored in `minio_data` volume

### 4. Redis (Rate Limiting)

- **Port**: 6379
- **Persistent**: Data stored in `redis_data` volume (AOF enabled)
- **Password**: Optional (set via `REDIS_PASSWORD` env)

## Environment Variables

Create a `.env` file in the project root:

```bash
# Application
APP_ENV=production
APP_PORT=8080

# Database
DB_USER=postgres
DB_PASS=your_secure_password
DB_NAME=go_rest_api_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-min-32-characters-long

# Storage
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin

# Redis (optional password)
REDIS_PASSWORD=

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RATE=100
RATE_LIMIT_AUTH_RATE=7
RATE_LIMIT_WINDOW_SECOND=60
```

## Production-Ready Features

### ✅ Automatic Redis Configuration

- Redis is **automatically enabled** in Docker
- No code changes needed
- Distributed rate limiting works out of the box

### ✅ Health Checks

All services have health checks:

- **PostgreSQL**: `pg_isready`
- **MinIO**: HTTP health endpoint
- **Redis**: `redis-cli ping`

### ✅ Service Dependencies

The API service waits for all dependencies to be healthy before starting:

- Database must be ready
- MinIO must be ready
- Redis must be ready

### ✅ Persistent Storage

All data is persisted in Docker volumes:

- `db_data`: PostgreSQL data
- `minio_data`: MinIO object storage
- `redis_data`: Redis data (with AOF persistence)

### ✅ Network Isolation

All services run on a private Docker network (`app_network`)

## Redis Configuration

Redis is configured with:

- **AOF Persistence**: Enabled for data durability
- **Password Protection**: Optional (set `REDIS_PASSWORD` env)
- **Health Checks**: Automatic monitoring
- **Volume Persistence**: Data survives container restarts

## Accessing Services

### From Host Machine

- API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- MinIO API: `localhost:9000`
- MinIO Console: `http://localhost:9001`
- Redis: `localhost:6379`

### From Other Containers

- API: `api:8080`
- PostgreSQL: `db:5432`
- MinIO: `minio:9000`
- Redis: `redis:6379`

## Troubleshooting

### Redis Connection Issues

```bash
# Check Redis logs
docker-compose logs redis

# Test Redis connection
docker-compose exec redis redis-cli ping

# If password is set
docker-compose exec redis redis-cli -a $REDIS_PASSWORD ping
```

### API Not Starting

```bash
# Check all service health
docker-compose ps

# View API logs
docker-compose logs api

# Check if Redis is accessible
docker-compose exec api ping redis
```

### Reset Everything

```bash
# Stop and remove all containers and volumes
docker-compose down -v

# Rebuild and start
docker-compose up -d --build
```

## Production Deployment

For production, consider:

1. **Set Strong Passwords**: Update all default passwords
2. **Use Secrets**: Don't commit `.env` file to git
3. **Enable Redis Password**: Set `REDIS_PASSWORD` for security
4. **Resource Limits**: Add CPU/memory limits to services
5. **Backup Strategy**: Regular backups of volumes
6. **Monitoring**: Add monitoring for all services
7. **SSL/TLS**: Use reverse proxy (nginx/traefik) for HTTPS

## Example Production docker-compose.yml Override

Create `docker-compose.prod.yml`:

```yaml
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
    restart: always

  redis:
    command: >
      redis-server --appendonly yes --requirepass ${REDIS_PASSWORD} --maxmemory
      256mb --maxmemory-policy allkeys-lru
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
    restart: always
```

Run with:

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```
