package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
	"github.com/rs/zerolog/log"
)

// RateLimiter implements a simple in-memory rate limiter using token bucket algorithm
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
	cleanup  *time.Ticker
}

type visitor struct {
	lastSeen time.Time
	count    int
	resetAt  time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed
// window: time window (e.g., 1 minute)
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
		cleanup:  time.NewTicker(5 * time.Minute), // cleanup old visitors every 5 minutes
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			v.mu.Lock()
			if now.Sub(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			lastSeen: time.Now(),
			resetAt:  time.Now().Add(rl.window),
		}
		rl.visitors[ip] = v
	}
	return v
}

func (rl *RateLimiter) allow(ip string) bool {
	v := rl.getVisitor(ip)

	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	v.lastSeen = now

	// Reset if window has passed
	if now.After(v.resetAt) {
		v.count = 0
		v.resetAt = now.Add(rl.window)
	}

	// Check if limit exceeded
	if v.count >= rl.rate {
		return false
	}

	v.count++
	return true
}

// RateLimit middleware limits requests per IP address
// Default: 100 requests per minute
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	if rate <= 0 {
		rate = 100 // default: 100 requests per minute
	}
	if window <= 0 {
		window = time.Minute
	}

	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			log.Ctx(c.Request.Context()).Warn().
				Str("ip", ip).
				Str("path", c.Request.URL.Path).
				Msg("rate limit exceeded")

			httpx.Fail(
				c,
				http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				"too many requests, please try again later",
				nil,
			)
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.rate))
		// Note: Remaining count could be enhanced with more sophisticated tracking

		c.Next()
	}
}
