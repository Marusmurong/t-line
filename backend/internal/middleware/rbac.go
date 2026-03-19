package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/pkg/response"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		if _, ok := roleSet[userRole]; !ok {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin", "super_admin")
}
