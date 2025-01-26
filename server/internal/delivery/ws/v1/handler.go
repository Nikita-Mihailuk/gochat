package v1

import (
	"github.com/Nikita-Mihailuk/gochat/server/internal/domain"
	"github.com/Nikita-Mihailuk/gochat/server/internal/service"
	"github.com/Nikita-Mihailuk/gochat/server/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type HandlerV1WS struct {
	services service.Services
	logger   *zap.Logger
	upgrader websocket.Upgrader
	rooms    map[uint][]domain.Client
	roomsMu  sync.Mutex
}

func NewHandlerV1WS(services service.Services) *HandlerV1WS {
	return &HandlerV1WS{
		services: services,
		logger:   logging.GetLogger(),
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		rooms:    make(map[uint][]domain.Client),
	}
}

func (h *HandlerV1WS) InitHandlerV1WS(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	h.RegisterWSRouts(v1)
}
