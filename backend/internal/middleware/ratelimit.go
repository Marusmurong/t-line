package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/t-line/backend/internal/pkg/response"
)

func RateLimit(rdb *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.FullPath())

		ctx := context.Background()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > int64(maxRequests) {
			response.Fail(c, http.StatusTooManyRequests, 40099, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
