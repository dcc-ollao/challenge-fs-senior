package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateBucket struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]rateBucket

	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]rateBucket),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) allow(key string) (allowed bool, remaining int, resetAt time.Time) {
	now := time.Now().UTC()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[key]
	if !ok || now.After(b.resetAt) {
		b = rateBucket{
			count:   0,
			resetAt: now.Add(rl.window),
		}
	}

	if b.count >= rl.limit {
		return false, 0, b.resetAt
	}

	b.count++
	rl.buckets[key] = b

	return true, rl.limit - b.count, b.resetAt
}

// RateLimit limits by JWT userId if present (set by AuthRequired middleware),
// otherwise falls back to IP (for public endpoints like /auth/login).
func RateLimit(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := ""

		if v, ok := c.Get(ContextUserIDKey); ok {
			if s, ok := v.(string); ok && s != "" {
				key = "user:" + s
			}
		}
		if key == "" {
			key = "ip:" + c.ClientIP()
		}

		allowed, remaining, resetAt := rl.allow(key)

		c.Header("X-RateLimit-Limit", itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", itoa(remaining))
		c.Header("X-RateLimit-Reset", itoa(int(resetAt.Unix())))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [32]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + (n % 10))
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
