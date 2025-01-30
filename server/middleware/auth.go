package middleware

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/service/token_manager"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenCookie, err := c.Request.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "access_token cookie отсутствует"})
			c.Abort()
			return
		}
		accessToken := accessTokenCookie.Value
		tokenManager := token_manager.GetTokenManager()
		claims, err := tokenManager.Parse(accessToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный или истекший access токен"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)

		c.Next()
	}
}
