package v1

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *HandlerV1HTTP) RegisterRoomsRouts(v1 *gin.RouterGroup) {
	rooms := v1.Group("/rooms", middleware.CheckAccessToken())

	rooms.GET("/", h.getAllRoomsHandler)
	rooms.POST("/", h.createRoomHandler)
	rooms.GET("/:id", h.getRoomMessagesHandler)
}

func (h *HandlerV1HTTP) getAllRoomsHandler(c *gin.Context) {
	rooms, err := h.services.Rooms.GetRoomsService()
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}
func (h *HandlerV1HTTP) createRoomHandler(c *gin.Context) {
	var input domain.InputRoomDTO
	if err := c.BindJSON(&input); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	err := h.services.Rooms.CreateRoomService(input)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Комната создана", zap.String("room_name", input.Name))
}

func (h *HandlerV1HTTP) getRoomMessagesHandler(c *gin.Context) {
	roomID := c.Param("id")
	messages, err := h.services.Rooms.GetRoomMessageService(roomID)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}
