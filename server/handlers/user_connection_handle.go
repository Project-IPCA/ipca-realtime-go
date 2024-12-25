package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	"github.com/Project-IPCA/ipca-realtime-go/repositories"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
)

type ResponseData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type userConnectionHandler struct {
	server *s.Server
}

func NewuserConnectionHandler(server *s.Server) *userConnectionHandler {
	return &userConnectionHandler{server: server}
}

func (userConnectionHanuserConnectionHandler *userConnectionHandler) ConsumeUserConnection(c echo.Context) error {
	subscriber := redis_client.Init(userConnectionHanuserConnectionHandler.server.Config)
	userRepository := repositories.NewUserRepository(userConnectionHanuserConnectionHandler.server.DB)
	defer subscriber.Close()

	userID := c.Param("user_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	userRepository.UpdateUserStatus(userID, true)

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("user-event:%s", userID))
	defer pubsub.Close()

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var data OnlineStudentMessage
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				fmt.Println(msg.Payload)
				fmt.Println("Failed to unmarshal message:", err)
				return
			}

			var sendData ResponseData
			if data.Action == "repeat-login" {
				sendData = ResponseData{
					Status:  false,
					Message: "repeat-login",
				}
			} else if data.Action == "reject-submission" {
				sendData = ResponseData{
					Status:  true,
					Message: "reject-submission",
				}
			} else if data.Action == "can-submit" {
				sendData = ResponseData{
					Status:  true,
					Message: "can-submit",
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

	userRepository.UpdateUserStatus(userID, false)
	pubsub.Unsubscribe(context.Background())
	pubsub.Close()

	return nil
}
