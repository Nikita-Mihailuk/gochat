package v1

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *HandlerV1WS) RegisterWSRouts(v1 *gin.RouterGroup) {
	ws := v1.Group("/ws")
	ws.GET("/:id", h.webSocketHandler)
	authorizationGroup := ws.Group("", middleware.CheckAccessToken())
	authorizationGroup.GET("/", h.getUserID)
}

func (h *HandlerV1WS) getUserID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *HandlerV1WS) webSocketHandler(c *gin.Context) {
	roomID, user, conn, err := h.initializeWebSocket(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer h.notifyParticipants(roomID)
	defer h.notificationUser(user.Name, roomID, "exit")
	defer h.cleanupConnection(roomID, user.ID)

	h.notifyParticipants(roomID)
	h.notificationUser(user.Name, roomID, "entrance")
	h.logger.Info("Пользователь присоединился к чату", zap.Uint("user_id", user.ID), zap.Uint("room_id", roomID))

	for {
		var msg domain.InputMessageDTO
		if err = conn.ReadJSON(&msg); err != nil {
			break
		}

		message := domain.InputMessageDTO{
			RoomID:  roomID,
			UserID:  msg.UserID,
			Content: msg.Content,
		}
		if err = h.services.Rooms.CreateMessageService(message); err != nil {
			h.logger.Error("Сообщение не сохранено", zap.Error(err))
			continue
		}

		tempMessage := domain.OutputMessageDTO{
			UserID:     user.ID,
			UserAvatar: user.PhotoURL,
			Content:    message.Content,
			UserName:   user.Name,
		}
		h.broadcastMessage(roomID, tempMessage)
	}
}

func (h *HandlerV1WS) initializeWebSocket(c *gin.Context) (uint, domain.User, *websocket.Conn, error) {
	roomID, _ := strconv.Atoi(c.Param("id"))
	userIDstr := c.Query("user_id")
	userID, _ := strconv.Atoi(userIDstr)
	user, err := h.services.User.GetProfileService(uint(userID))
	if err != nil {
		return 0, domain.User{}, nil, fmt.Errorf("Пользователь с ID %d не найден", userID)
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Ошибка обновления WebSocket: " + err.Error())
		return 0, domain.User{}, nil, fmt.Errorf("не удалось установить WebSocket соединение")
	}

	client := domain.Client{
		User: user,
		Conn: conn,
	}

	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()
	h.rooms[uint(roomID)] = append(h.rooms[uint(roomID)], client)
	return uint(roomID), user, conn, nil
}

func (h *HandlerV1WS) cleanupConnection(roomID, userID uint) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	clients := h.rooms[roomID]
	for i, client := range clients {
		if client.User.ID == userID {
			h.rooms[roomID] = append(clients[:i], clients[i+1:]...)
			h.logger.Info("Пользователь вышел из чата", zap.Uint("user_id", userID), zap.Uint("room_id", roomID))
			return
		}
	}
}

func (h *HandlerV1WS) broadcastMessage(roomID uint, message domain.OutputMessageDTO) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	for _, client := range h.rooms[roomID] {
		if err := client.Conn.WriteJSON(gin.H{"type": "message", "message": message}); err != nil {
			h.logger.Error("Ошибка отправки сообщения", zap.Error(err))
		}
	}
}

func (h *HandlerV1WS) notificationUser(userName string, roomID uint, style string) {
	var message string
	if style == "entrance" {
		message = "Пользователь с именем " + userName + " присоединился к чату"
	} else if style == "exit" {
		message = "Пользователь с именем " + userName + " вышел из чата"
	} else {
		h.logger.Error("Неправильный тип уведомления")
		return
	}

	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	for _, client := range h.rooms[roomID] {
		if err := client.Conn.WriteJSON(gin.H{
			"type":    "notification",
			"message": message,
		}); err != nil {
			h.logger.Error("Ошибка отправки уведомления", zap.Error(err))
		}
	}
}

func (h *HandlerV1WS) notifyParticipants(roomID uint) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	users := make([]domain.User, 0, len(h.rooms[roomID]))
	for _, client := range h.rooms[roomID] {
		users = append(users, client.User)
	}

	for _, client := range h.rooms[roomID] {
		if err := client.Conn.WriteJSON(gin.H{
			"type":         "participants",
			"participants": users,
		}); err != nil {
			h.logger.Error("Ошибка обновления списка участников", zap.Error(err))
		}
	}
}
