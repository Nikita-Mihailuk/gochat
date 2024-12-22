package users

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/cfg"
	"github.com/Nikita-Mihailuk/gochat/server/internal/users/model"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	db     *gorm.DB
	logger *zap.Logger
}

func GetHandler(db *gorm.DB) *Handler {
	return &Handler{
		db:     db,
		logger: logging.GetLogger(),
	}
}

func (h *Handler) RegisterRouter(router *gin.Engine) {
	router.POST("/register", h.RegisterHandler)
	router.POST("/login", h.LoginHandler)
	router.GET("/rooms", h.ListRoomsHandler)
	router.POST("/rooms", h.CreateRoomHandler)
	router.GET("/rooms/:id", h.RoomMessagesHandler)
	router.GET("/ws/:roomID", h.WebSocketHandler)
	router.GET("/profile/:id", h.ProfileHandler)
	router.PATCH("/profile/:id", h.ProfileUpdate)
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := c.BindJSON(&data); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	if len(data.Email) > 30 || len(data.Name) > 30 || len(data.Password) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Превышено количество допустимых символов"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	user := model.User{Email: data.Email, PasswordHash: string(hashedPassword), Name: data.Name}
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с такой почтой уже существует"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Регистрация прошла успешно"})
	h.logger.Info("Пользователь с почтой " + user.Email + " зарегистрировался")
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&data); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	var user model.User
	if err := h.db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вход выполнен успешно", "user_name": user.Name, "user_id": user.ID})
	h.logger.Info("Пользователь с ID " + strconv.Itoa(int(user.ID)) + " вошёл в систему")
}

func (h *Handler) ListRoomsHandler(c *gin.Context) {
	var roomsList []model.Room
	h.db.Find(&roomsList)
	c.JSON(http.StatusOK, roomsList)
}

func (h *Handler) CreateRoomHandler(c *gin.Context) {
	var data struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&data); err != nil {
		h.logger.Error("Недействительный запрос")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный запрос"})
		return
	}
	if len(data.Name) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Превышено допустимое количество символов"})
		return
	}
	room := model.Room{Name: data.Name}
	if err := h.db.Create(&room).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Комната с таким именем уже существует"})
		return
	}
	c.JSON(http.StatusCreated, room)
	h.logger.Info("Комната с именем " + room.Name + " создана")
}

func (h *Handler) RoomMessagesHandler(c *gin.Context) {
	roomID := c.Param("id")
	var messages []model.Message

	if err := h.db.Table("messages").
		Select("messages.*, users.photo_url AS user_avatar").
		Joins("JOIN users ON users.id = messages.user_id").
		Where("messages.room_id = ?", roomID).
		Order("messages.created_at").
		Find(&messages).Error; err != nil {
		h.logger.Warn("Комната не найдена")
		c.JSON(http.StatusNotFound, gin.H{"error": "Комната не найдена"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	rooms      = make(map[uint][]*websocket.Conn)
	roomsUsers = make(map[uint]map[uint]model.User)
	roomsMu    sync.Mutex
)

func (h *Handler) WebSocketHandler(c *gin.Context) {

	roomIDint, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		h.logger.Error("Недействительный ID комнаты %s" + c.Param("roomID"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный ID комнаты"})
		return
	}
	roomID := uint(roomIDint)

	var room model.Room
	if err = h.db.First(&room, roomID).Error; err != nil {
		h.logger.Error("Комната с таким ID не найдена: " + strconv.Itoa(int(roomID)))
		c.JSON(http.StatusNotFound, gin.H{"error": "Комната не найдена"})
		return
	}

	userIDint, err := strconv.Atoi(c.Query("user_id"))

	if err != nil {
		h.logger.Error("Недействительный ID пользователя: " + c.Query("user_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный ID пользователя"})
		return
	}
	userID := uint(userIDint)

	var user model.User
	if err = h.db.First(&user, userID).Error; err != nil {
		h.logger.Error("Пользователь с таким ID не найден: " + strconv.Itoa(int(userID)))
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Ошибка обновления WebSocket: " + err.Error())
		return
	}
	defer conn.Close()

	roomsMu.Lock()
	if roomsUsers[roomID] == nil {
		roomsUsers[roomID] = make(map[uint]model.User)
	}
	roomsUsers[roomID][user.ID] = user
	rooms[roomID] = append(rooms[roomID], conn)
	roomsMu.Unlock()

	notifyParticipants(roomID, h.logger)
	notificationUser(user.Name, roomID, "entrance", h.logger)
	h.logger.Info("Пользователь с именем " + user.Name + " присоединился к чату " + strconv.Itoa(int(roomID)))

	defer func() {
		roomsMu.Lock()
		delete(roomsUsers[roomID], user.ID)
		for i, client := range rooms[roomID] {
			if client == conn {
				rooms[roomID] = append(rooms[roomID][:i], rooms[roomID][i+1:]...)
				break
			}
		}
		roomsMu.Unlock()

		notifyParticipants(roomID, h.logger)
		notificationUser(user.Name, roomID, "exit", h.logger)
		h.logger.Info("Пользователь с именем " + user.Name + " вышел из чата " + strconv.Itoa(int(roomID)))
	}()

	for {
		var msg struct {
			UserID  uint   `json:"user_id"`
			Message string `json:"message"`
		}

		if err = conn.ReadJSON(&msg); err != nil {
			h.logger.Error("Ошибка чтения сообщения WebSocket: " + err.Error())
			break
		}

		message := model.Message{
			RoomID:   room.ID,
			UserID:   msg.UserID,
			UserName: user.Name,
			Content:  msg.Message,
		}

		if err = h.db.Create(&message).Error; err != nil {
			h.logger.Error("Сообщение не сохранено: " + err.Error())
			continue
		}

		tempMessage := model.Message{
			UserID:     user.ID,
			UserAvatar: user.PhotoURL,
			Content:    message.Content,
			UserName:   user.Name,
		}

		roomsMu.Lock()

		for _, client := range rooms[room.ID] {
			if err = client.WriteJSON(gin.H{
				"type":    "message",
				"message": tempMessage,
			}); err != nil {
				h.logger.Error("Ошибка отправки сообщения: " + err.Error())
			}
		}
		roomsMu.Unlock()
	}
}

func notificationUser(userName string, roomID uint, style string, logger *zap.Logger) {
	var message string
	if style == "entrance" {
		message = "Пользователь с именем " + userName + " присоединился к чату"
	} else if style == "exit" {
		message = "Пользователь с именем " + userName + " вышел из чата"
	} else {
		logger.Error("Неправильный тип уведомления")
		return
	}
	roomsMu.Lock()
	for _, client := range rooms[roomID] {
		if err := client.WriteJSON(gin.H{
			"type":    "notification",
			"message": message,
		}); err != nil {
			logger.Error("Ошибка отправки уведомления: " + err.Error())
		}
	}
	roomsMu.Unlock()
}

func notifyParticipants(roomID uint, logger *zap.Logger) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	users := make([]model.User, 0, len(roomsUsers[roomID]))
	for _, user := range roomsUsers[roomID] {
		users = append(users, user)
	}

	for _, conn := range rooms[roomID] {
		if err := conn.WriteJSON(gin.H{
			"type":         "participants",
			"participants": users,
		}); err != nil {
			logger.Error("Ошибка обновления списка участников: " + err.Error())
		}
	}
}

func (h *Handler) ProfileHandler(c *gin.Context) {
	userID := c.Param("id")
	var user model.User
	if err := h.db.Where("id = ?", userID).Find(&user).Error; err != nil {
		h.logger.Error("Пользователь с таким ID не найден: %" + userID)
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) ProfileUpdate(c *gin.Context) {
	userID := c.Param("id")

	currentPassword := c.PostForm("current_password")
	if currentPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите текущий пароль"})
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный текущий пароль"})
		return
	}

	newName := c.PostForm("name")
	newPassword := c.PostForm("new_password")

	if newName != "" {
		user.Name = newName
	}

	if newPassword != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		user.PasswordHash = string(hashedPassword)
	}

	file, err := c.FormFile("photo")
	if err == nil {
		filePath := fmt.Sprintf("uploads/%d_%s", user.ID, file.Filename)
		if err = c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения фото"})
			return
		}
		getConfig := cfg.GetConfig()
		user.PhotoURL = fmt.Sprintf("http://%s:%s/%s", getConfig.Listen.BindIP, getConfig.Listen.Port, filePath)
	}

	user.UpdatedAt = time.Now()

	if err = h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Профиль успешно обновлён",
		"user":    user,
	})
	h.logger.Info("Пользователь с именем " + user.Name + " обновил профиль")
}
