// Package middleware 提供 HTTP 中间件
package middleware

import (
	"net/http"
	"strings"

	"go-cargo/internal/config"
	"go-cargo/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 支持 "Bearer <token>" 格式
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 解析 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.Global.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "认证令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// 提取 Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.Response{
				Code:    401,
				Message: "认证令牌解析失败",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// AdminOnly 管理员权限中间件
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			c.JSON(http.StatusForbidden, models.Response{
				Code:    403,
				Message: "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}
