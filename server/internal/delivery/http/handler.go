package http

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/delivery/http/v1"
	"github.com/Nikita-Mihailuk/gochat/server/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services service.Services
}

func NewHandlerHTTP(services service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitHandlerHTTP(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	handlerV1 := v1.NewHandlerV1HTTP(h.services)
	api := router.Group("/api")
	handlerV1.InitHandlerV1HTTP(api)
}
