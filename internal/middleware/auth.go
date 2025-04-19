package middleware

import (
	"net/http"
	"strings"
	"video-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在测试模式下，如果已经设置了userId，则跳过token验证
		if gin.Mode() == gin.TestMode {
			if _, exists := c.Get("userId"); exists {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "未授权",
				"data": nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "无效的认证格式",
				"data": nil,
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "无效的token",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func SetUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在测试模式下，如果已经设置了userId，则跳过token验证
		if gin.Mode() == gin.TestMode {
			if _, exists := c.Get("userId"); exists {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Next()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
