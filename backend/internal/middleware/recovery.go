package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/pkg/logger"
	"github.com/t-line/backend/internal/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L.Errorf("panic recovered: %v", r)
				response.Fail(c, http.StatusInternalServerError, 50000, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}
