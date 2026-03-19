package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/t-line/backend/internal/pkg/logger"
)

const RequestIDKey = "X-Request-ID"

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDKey)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Set(RequestIDKey, requestID)
		c.Header(RequestIDKey, requestID)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logger.L.Infow("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"request_id", requestID,
		)
	}
}
