package v1

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *HandlerV1HTTP) RegisterUserRouts(v1 *gin.RouterGroup) {
	users := v1.Group("/users")
	users.POST("/register", h.registerUserHandler)
	users.POST("/login", h.loginUserHandler)
	users.GET("/:id", h.getProfileUserHandler)
	users.PATCH("/:id", h.updateProfileUserHandler)
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
	userID, err := h.services.User.LoginUserService(input)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
	h.logger.Info("Пользователь вошёл в систему", zap.Uint("user_id", userID))
}

func (h *HandlerV1HTTP) getProfileUserHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	user, err := h.services.User.GetProfileService(uint(userID))
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1HTTP) updateProfileUserHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	file, _ := c.FormFile("photo")
	var updateUser = domain.UpdateUserDTO{
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
	c.JSON(http.StatusOK, user)
	h.logger.Info("Пользователь обновил профиль", zap.Uint("user_id", user.ID))
}
