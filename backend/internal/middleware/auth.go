package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/t-line/backend/internal/pkg/jwt"
	"github.com/t-line/backend/internal/pkg/response"
)

const (
	CtxUserID      = "user_id"
	CtxUserRole    = "user_role"
	CtxMemberLevel = "member_level"
)

func JWTAuth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证信息")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := jwtMgr.ParseToken(token)
		if err != nil {
			response.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		// Only accept access tokens for API requests
		if claims.TokenType != jwt.TokenTypeAccess {
			response.Unauthorized(c, "无效的令牌类型")
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUserRole, claims.Role)
		c.Set(CtxMemberLevel, claims.MemberLevel)
		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	v, _ := c.Get(CtxUserID)
	id, _ := v.(int64)
	return id
}

func GetUserRole(c *gin.Context) string {
	v, _ := c.Get(CtxUserRole)
	role, _ := v.(string)
	return role
}
