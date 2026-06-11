package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecretForAuth())

func getJWTSecretForAuth() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "fukubukuro-secret-key"
}

// AuthMiddleware JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效或已过期"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token解析失败"})
			c.Abort()
			return
		}

		// 把用户信息存入上下文，方便后续接口使用
		c.Set("user_id", fmt.Sprintf("%v", claims["user_id"]))
		c.Set("username", claims["username"])
		c.Next()
	}
}
