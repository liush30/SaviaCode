package middleware

import (
	"eldercare_health/app/internal/pkg/tool"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware 认证中间件，检查 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头缺少Authorization"})
			c.Abort()
			return
		}

		// 提取 Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // 检查是否有 Bearer 前缀
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式不正确"})
			c.Abort()
			return
		}

		// 解析并验证 token
		claims, err := tool.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "不合法的token"})
			c.Abort()
			return
		}

		// 将解析出来的用户名保存到上下文中，后续可以用于其他操作
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
