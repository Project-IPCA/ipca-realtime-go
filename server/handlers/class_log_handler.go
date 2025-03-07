package handlers

import (
	"context"
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/Project-IPCA/ipca-realtime-go/models"
	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	"github.com/Project-IPCA/ipca-realtime-go/repositories"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/labstack/echo/v4"
)

type ClassLogHandler struct {
	server *s.Server
}

func NewClassLogHandler(server *s.Server) *ClassLogHandler {
	return &ClassLogHandler{server: server}
}

func (handler *ClassLogHandler) ConsumeClassLog(c echo.Context) error {
	subscriber := redis_client.Init(handler.server.Config)
	classLogRepository := repositories.NewClassLogRepository(handler.server.DB)
	defer subscriber.Close()

	groupID := c.Param("group_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	var classLog []models.ActivityLog
	classLogRepository.GetActivityLog(&classLog, groupID)

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("logs:%s", groupID))
	defer pubsub.Close()

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg.Payload)
			flusher.Flush()
		}
	}()

	<-c.Request().Context().Done()

	pubsub.Unsubscribe(context.Background())
	pubsub.Close()

	return nil
}
