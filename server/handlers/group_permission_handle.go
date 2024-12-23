package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
)

type GroupPermissionHandler struct {
	server *s.Server
}

func NewGroupPermissionHandler(server *s.Server) *GroupPermissionHandler {
	return &GroupPermissionHandler{server: server}
}

func (GroupPermissionHandler *GroupPermissionHandler) ConsumeGroupPermission(c echo.Context) error {
	subscriber := redis_client.Init(GroupPermissionHandler.server.Config)
	defer subscriber.Close()

	groupID := c.Param("group_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("group-permission:%s", groupID))
	defer pubsub.Close()

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var message string
			var sendData ResponseData

			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				fmt.Println(msg.Payload)
				fmt.Println("Failed to unmarshal message:", err)
				continue
			}

			if message == "permission-change" {
				sendData = ResponseData{
					Status:  true,
					Message: message,
				}
			} else if message == "logout-all" {
				sendData = ResponseData{
					Status:  false,
					Message: message,
				}
			} else {
				sendData = ResponseData{
					Status:  false,
					Message: "something-wrong",
				}
			}
			sendDataRaw, err := json.Marshal(sendData)
			if err != nil {
				fmt.Println("Failed to marshal updated data:", err)
				return
			}
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", sendDataRaw)
			flusher.Flush()
		}
	}()

	<-c.Request().Context().Done()

	pubsub.Unsubscribe(context.Background())
	pubsub.Close()

	return nil
}
