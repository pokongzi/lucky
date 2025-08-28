package middleware

import (
	"net/http"
	"strings"

	"lucky/common/mysql"
	"lucky/model"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少认证信息",
			})
			c.Abort()
			return
		}

		// 提取Bearer token
		token := extractBearerToken(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证格式",
			})
			c.Abort()
			return
		}

		// 验证token
		authService := service.NewAuthService(mysql.DB)
		user, err := authService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("open_id", user.OpenID)

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（即使没有token也可以继续）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 提取Bearer token
		token := extractBearerToken(authHeader)
		if token == "" {
			c.Next()
			return
		}

		// 验证token
		authService := service.NewAuthService(mysql.DB)
		user, err := authService.ValidateAccessToken(token)
		if err != nil {
			// 认证失败，但不阻止请求继续
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("open_id", user.OpenID)

		c.Next()
	}
}

// extractBearerToken 从Authorization header提取Bearer token
func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetCurrentUser 从gin.Context获取当前用户
func GetCurrentUser(c *gin.Context) (*model.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userModel, ok := user.(*model.User)
	return userModel, ok
}

// GetCurrentUserID 从gin.Context获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint64)
	return id, ok
}
