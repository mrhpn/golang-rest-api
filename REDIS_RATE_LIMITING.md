# Redis-Based Rate Limiting Implementation

## Overview

This implementation addresses all three critical issues with in-memory rate
limiting:

1. ✅ **Multi-Instance Isolation**: Shared Redis state across all replicas
2. ✅ **Fixed Window Bursting**: Sliding window algorithm prevents edge attacks
3. ✅ **X-RateLimit-Reset Header**: Proper Unix timestamp for client retry logic

## Implementation Details

### Algorithm: Sliding Window Log

The implementation uses a **sliding window log** algorithm stored in Redis:

- Each request is stored as a sorted set entry with timestamp as score
- Old entries (outside the window) are automatically removed
- Count of entries in the window determines if request is allowed
- Reset time is calculated from the oldest entry in the window

### Key Features

1. **Distributed Rate Limiting**: All API instances share the same Redis state
2. **Sliding Window**: Prevents burst attacks at window boundaries
3. **Atomic Operations**: Uses Redis pipelines for consistency
4. **Automatic Cleanup**: Keys expire automatically after window duration
5. **Proper Headers**: Includes all RFC 7231 compliant rate limit headers

## Configuration

### Environment Variables

```bash
# Enable Redis (required for distributed rate limiting)
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=          # Optional, leave empty if no password
REDIS_DB=0               # Redis database number (0-15)

# Rate limiting configuration
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RATE=100      # Requests per window
RATE_LIMIT_WINDOW_SECOND=60  # Window in seconds
```

### Fallback Behavior

- If Redis is disabled or unavailable, the system falls back to in-memory rate
  limiting
- A warning is logged when falling back (not suitable for multi-instance
  deployments)

## HTTP Headers

The middleware sets the following headers (RFC 7231 compliant):

- `X-RateLimit-Limit`: Maximum number of requests allowed
- `X-RateLimit-Remaining`: Number of requests remaining in current window
- `X-RateLimit-Reset`: Unix timestamp when the rate limit resets

### Example Response Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067200
```

## Production Costs & Considerations

### Redis Costs

#### Cloud Provider Options

1. **AWS ElastiCache (Redis)**

   - **t3.micro**: ~$15/month (development)
   - **t3.small**: ~$30/month (small production)
   - **t3.medium**: ~$60/month (medium production)
   - **Memory optimized**: $100-500+/month (high traffic)

2. **Google Cloud Memorystore**

   - **Basic tier**: ~$30-100/month
   - **Standard tier**: ~$100-500+/month

3. **Azure Cache for Redis**

   - **Basic C0**: ~$15/month
   - **Standard C1**: ~$60/month
   - **Premium**: $200+/month

4. **Self-Hosted (EC2/Docker)**
   - **t3.micro instance**: ~$7-10/month
   - **t3.small instance**: ~$15-20/month
   - Requires maintenance and monitoring

#### Cost Optimization Tips

1. **Use Redis for Multiple Purposes**:

   - Rate limiting
   - Session storage
   - Caching
   - Job queues
   - This spreads the cost across features

2. **Start Small**:

   - Begin with smallest instance
   - Monitor memory usage
   - Scale up only when needed

3. **Use Redis Cluster for High Availability**:

   - Only if you need 99.99% uptime
   - Adds ~2x cost

4. **Consider Managed Services**:
   - Less operational overhead
   - Automatic backups and updates
   - Better for small teams

### Memory Usage

For rate limiting alone:

- **Per IP**: ~100-200 bytes per active IP
- **10,000 active IPs**: ~1-2 MB
- **100,000 active IPs**: ~10-20 MB

Redis is very memory-efficient for this use case.

### Network Costs

- **Minimal**: Each rate limit check = 1-2 Redis commands
- **Latency**: < 1ms for local Redis, < 5ms for cloud Redis
- **Bandwidth**: Negligible (few KB per request)

### Operational Considerations

#### Benefits

1. ✅ **Shared State**: Works across multiple instances
2. ✅ **Scalability**: Can handle millions of requests
3. ✅ **Persistence**: Optional persistence for audit trails
4. ✅ **Monitoring**: Redis provides built-in metrics
5. ✅ **Flexibility**: Can implement complex rate limiting rules

#### Drawbacks

1. ⚠️ **Additional Infrastructure**: One more service to manage
2. ⚠️ **Cost**: Additional monthly expense
3. ⚠️ **Latency**: Small network overhead (usually < 1ms)
4. ⚠️ **Dependency**: API depends on Redis availability

#### Mitigation Strategies

1. **Fail-Open**: If Redis fails, allow requests (current implementation)

   - Alternative: Fail-closed (reject all requests if Redis down)
   - Choose based on your security requirements

2. **Connection Pooling**: Already implemented (10 connections)
3. **Health Checks**: Monitor Redis in readiness probe
4. **Backup Strategy**: Use Redis persistence for critical data

## Comparison: In-Memory vs Redis

| Feature            | In-Memory       | Redis           |
| ------------------ | --------------- | --------------- |
| **Multi-Instance** | ❌ No           | ✅ Yes          |
| **Sliding Window** | ❌ Fixed window | ✅ True sliding |
| **Reset Header**   | ❌ Missing      | ✅ Included     |
| **Cost**           | Free            | $15-100+/month  |
| **Latency**        | < 0.1ms         | < 1-5ms         |
| **Scalability**    | Limited         | Unlimited       |
| **Persistence**    | ❌ No           | ✅ Optional     |

## When to Use Redis

### ✅ Use Redis When:

- Deploying multiple API instances (ECS, Kubernetes, etc.)
- Need accurate rate limiting across instances
- Want to prevent burst attacks
- Need rate limit metrics/analytics
- Already using Redis for other features

### ❌ Skip Redis When:

- Single instance deployment
- Very low traffic (< 1000 req/min)
- Budget constraints
- Simple use case (can accept fixed window)

## Migration Path

1. **Phase 1**: Deploy with Redis disabled (in-memory)
2. **Phase 2**: Enable Redis in staging
3. **Phase 3**: Monitor and tune
4. **Phase 4**: Enable in production

## Monitoring

### Key Metrics to Monitor

1. **Redis Memory Usage**: Should stay under 80% of instance size
2. **Redis CPU**: Should be < 50% for rate limiting
3. **Rate Limit Hits**: Track how many requests are rate limited
4. **Redis Latency**: Should be < 5ms p99
5. **Connection Pool**: Monitor active connections

### Health Checks

The readiness probe should check Redis connectivity:

```go
// Add to health check
if ctx.Redis != nil {
    err := ctx.Redis.Ping(ctx).Err()
    // Handle error
}
```

## Example: Cost Calculation

**Scenario**: 3 API instances, 10,000 requests/minute, 1,000 unique IPs

- **Redis Instance**: t3.small ($30/month)
- **Memory Usage**: ~200 KB (negligible)
- **Network**: < 1 GB/month (negligible)
- **Total Cost**: ~$30/month

**ROI**: Prevents DDoS attacks, ensures fair usage, enables scaling

## Conclusion

Redis-based rate limiting is **essential for production multi-instance
deployments**. The cost (~$15-100/month) is minimal compared to the benefits:

- ✅ Prevents bypassing rate limits across instances
- ✅ Prevents burst attacks
- ✅ Provides proper client feedback
- ✅ Enables horizontal scaling

For single-instance deployments or very low traffic, in-memory rate limiting may
be sufficient, but Redis is recommended for any production API.
