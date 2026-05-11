package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"doctor-go/internal/infrastructure/redis"
	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
)

func RateLimit(cache *redis.Client, action string, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cache == nil || limit <= 0 || window <= 0 {
			c.Next()
			return
		}

		key := "rate_limit:" + action + ":" + c.ClientIP()
		count, err := cache.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			_ = cache.Expire(c.Request.Context(), key, window).Err()
		}
		if count > limit {
			ttl := cache.TTL(c.Request.Context(), key).Val()
			if ttl > 0 {
				c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			}
			response.Fail(c, http.StatusTooManyRequests, appErrors.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
