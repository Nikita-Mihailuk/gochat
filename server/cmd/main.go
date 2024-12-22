package main

import (
	"fmt"
	"github.com/Nikita-Mihailuk/gochat/server/internal/cfg"
	"github.com/Nikita-Mihailuk/gochat/server/internal/users"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/dbClients"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logging.GetLogger()

	db := dbClients.GetDB()
	handler := users.GetHandler(db)

	logger.Info("Создание роутера")
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Static("/uploads", "./uploads")

	logger.Info("Регистрация эндпоинтов")
	handler.RegisterRouter(router)

	start(router)
}
func start(router *gin.Engine) {

	getConfig := cfg.GetConfig()
	logger := logging.GetLogger()

	logger.Info("Прослушивание сервера",
		zap.String("bind_ip", getConfig.Listen.BindIP),
		zap.String("port", getConfig.Listen.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", getConfig.Listen.BindIP, getConfig.Listen.Port))

	if err != nil {
		logger.Fatal(err.Error())
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Minute,
		ReadTimeout:  15 * time.Minute,
	}

	logger.Fatal(server.Serve(listener).Error())
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
