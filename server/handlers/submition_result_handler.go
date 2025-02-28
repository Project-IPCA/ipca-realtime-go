package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
)

type SubmitionResultHandler struct {
	server *s.Server
}

func NewSubmitionResultHandler(server *s.Server) *SubmitionResultHandler {
	return &SubmitionResultHandler{server: server}
}

func (submitionResultHandler *SubmitionResultHandler) ConsumeSubmitionResult(c echo.Context) error {
	subscriber := redis_client.Init(submitionResultHandler.server.Config)
	defer subscriber.Close()

	jobID := c.Param("job_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("submission-result:%s", jobID))
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

	<-c.Request().Context().Done()

	pubsub.Unsubscribe(context.Background())
	pubsub.Close()

	return nil
}
