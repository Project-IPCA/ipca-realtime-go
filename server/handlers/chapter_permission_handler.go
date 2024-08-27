package handlers

import (
	"fmt"
	"net/http"
	"context"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
)

type ChapterPermissionHandler struct {
	server *s.Server
}

func NewChapterPermissionHandler(server *s.Server) *ChapterPermissionHandler {
	return &ChapterPermissionHandler{server: server}
}

func (chapterPermissionHandler *ChapterPermissionHandler)ConsumeChapterPermission(c echo.Context) error {
	subscriber := redis_client.Init(chapterPermissionHandler.server.Config);
	defer subscriber.Close()

	groupID := c.Param("group_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
    c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
    c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("chapter-permission:%s", groupID))
	defer pubsub.Close()

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", msg.Payload)
			flusher.Flush()
		}
	}()

	<- c.Request().Context().Done()

	pubsub.Unsubscribe(context.Background())
	pubsub.Close()

	return nil
}