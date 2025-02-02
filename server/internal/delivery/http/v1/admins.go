package v1

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *HandlerV1HTTP) RegisterAdminRouts(v1 *gin.RouterGroup) {
	isAdminGroup := v1.Group("admins", middleware.IsAdmin())

	isAdminGroup.GET("/rooms", h.getAllRoomsHandler)
	isAdminGroup.PATCH("/rooms/:id", h.updateRoomHandler)
	isAdminGroup.DELETE("/rooms/:id", h.deleteRoomHandler)

	isAdminGroup.GET("/users", h.getAllUsersHandler)
	isAdminGroup.PATCH("/users/:id", h.updateUserHandler)
	isAdminGroup.DELETE("/users/:id", h.deleteUserHandler)

	isAdminGroup.GET("/sessions", h.getAllSessionsHandler)
	isAdminGroup.DELETE("/sessions/:id", h.DeleteSessionHandler)
}

func (h *HandlerV1HTTP) updateRoomHandler(c *gin.Context) {
	var input domain.InputRoomDTO
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	roomIDstr := c.Param("id")
	roomID, _ := strconv.Atoi(roomIDstr)
	err := h.services.Admin.UpdateRoomService(domain.Room{ID: uint(roomID), Name: input.Name})
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Изменена комната", zap.Int("room_id", roomID))
}

func (h *HandlerV1HTTP) deleteRoomHandler(c *gin.Context) {
	roomID := c.Param("id")
	err := h.services.Admin.DeleteRoomService(roomID)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Удалена комната", zap.String("room_id", roomID))
}

func (h *HandlerV1HTTP) getAllUsersHandler(c *gin.Context) {
	users, err := h.services.Admin.GetUsersService()
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1HTTP) updateUserHandler(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))

	file, _ := c.FormFile("photo")
	var updateUser = domain.UpdateUserDTO{
		UserId:     uint(userID),
		NewName:    c.PostForm("name"),
		FileHeader: file,
	}
	err := h.services.Admin.UpdateUserService(updateUser)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Пользователь обновлён", zap.Int("user_id", userID))
}

func (h *HandlerV1HTTP) deleteUserHandler(c *gin.Context) {
	userID := c.Param("id")
	err := h.services.Admin.DeleteUserService(userID)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Пользователь удалён", zap.String("user_id", userID))
}

func (h *HandlerV1HTTP) getAllSessionsHandler(c *gin.Context) {
	sessions, err := h.services.Admin.GetSessionsService()
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *HandlerV1HTTP) DeleteSessionHandler(c *gin.Context) {
	sessionID := c.Param("id")
	err := h.services.Admin.DeleteSessionService(sessionID)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Сессия удалена", zap.String("session_id", sessionID))
}
