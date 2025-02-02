package middleware

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/service/token_manager"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenCookie, err := c.Request.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access_token cookie отсутствует"})
			c.Abort()
			return
		}
		accessToken := accessTokenCookie.Value
		tokenManager := token_manager.GetTokenManager()
		claims, err := tokenManager.Parse(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или истекший access токен"})
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Нет доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}
