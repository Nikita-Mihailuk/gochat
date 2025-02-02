package v1

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *HandlerV1HTTP) RegisterUserRouts(v1 *gin.RouterGroup) {

	v1.POST("/register", h.registerUserHandler)
	v1.POST("/login", h.loginUserHandler)
	v1.POST("/auth/refresh", h.refreshTokens)

	authorizationGroup := v1.Group("/users", middleware.CheckAccessToken())
	authorizationGroup.GET("/", h.getProfileUserHandler)
	authorizationGroup.PATCH("/", h.updateProfileUserHandler)
	authorizationGroup.DELETE("/logout", h.deleteUserSession)
	authorizationGroup.GET("/auth/check", h.authCheck)
}

func (h *HandlerV1HTTP) registerUserHandler(c *gin.Context) {
	var input domain.InputUserDTO
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	err := h.services.User.RegisterUserService(input)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Пользователь зарегистрировался", zap.String("email", input.Email))
}

func (h *HandlerV1HTTP) loginUserHandler(c *gin.Context) {
	var input domain.InputUserDTO
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	tokens, err := h.services.User.LoginUserService(input)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token",
		tokens.RefreshToken,
		30*24*3600, // 30 дней
		"/",
		"",
		false, // secure true при использовании HTTPS
		true)

	c.SetCookie("access_token",
		tokens.AccessToken,
		900, // 15 минут
		"/",
		"",
		false, // secure true при использовании HTTPS
		true)

	h.logger.Debug("Пользователь с данными токенами вошел в систему",
		zap.String("accessToken", tokens.AccessToken),
		zap.String("refreshToken", tokens.RefreshToken))
}

func (h *HandlerV1HTTP) getProfileUserHandler(c *gin.Context) {
	userIDstr, _ := c.Get("user_id")
	userID, _ := strconv.Atoi(fmt.Sprint(userIDstr))
	user, err := h.services.User.GetProfileService(uint(userID))
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1HTTP) updateProfileUserHandler(c *gin.Context) {
	userIDstr, _ := c.Get("user_id")
	userID, _ := strconv.Atoi(fmt.Sprint(userIDstr))
	file, _ := c.FormFile("photo")
	var updateUser = domain.UpdateProfileDTO{
		UserId:          uint(userID),
		CurrentPassword: c.PostForm("current_password"),
		NewPassword:     c.PostForm("new_password"),
		NewName:         c.PostForm("name"),
		FileHeader:      file,
	}

	user, err := h.services.User.UpdateProfileService(updateUser)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Пользователь обновил профиль", zap.Uint("user_id", user.ID))
	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1HTTP) refreshTokens(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		h.logger.Error("Refresh токен не найден")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh токен не найден"})
		return
	}

	newTokens, err := h.services.User.RefreshTokens(refreshToken)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token",
		newTokens.AccessToken,
		900, /// 15 минут
		"",
		"",
		false, // secure true при использовании HTTPS
		true)

	c.SetCookie("refresh_token",
		newTokens.RefreshToken,
		30*24*360, // 30 дней
		"",
		"",
		false, // secure true при использовании HTTPS
		true)
}

func (h *HandlerV1HTTP) deleteUserSession(c *gin.Context) {
	userID, _ := c.Get("user_id")
	err := h.services.User.DeleteSessionServiceByUserID(fmt.Sprint(userID))
	if err != nil {
		h.logger.Error("Ошибка при удалении сессии пользователя", zap.String("user_id", fmt.Sprint(userID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при удалении сессии пользователя"})
		return
	}

	c.SetCookie("access_token",
		"",
		-1,
		"",
		"",
		false, // secure true при использовании HTTPS
		true)

	c.SetCookie("refresh_token",
		"",
		-1,
		"",
		"",
		false, // secure true при использовании HTTPS
		true)

	h.logger.Info("Пользователь вышел из аккаунта", zap.String("user_id", fmt.Sprint(userID)))
}

func (h *HandlerV1HTTP) authCheck(c *gin.Context) {
	role, err := c.Get("role")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Роль не передана"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"role": role})
}
