package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Project-IPCA/ipca-realtime-go/redis_client"
	"github.com/Project-IPCA/ipca-realtime-go/repositories"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/labstack/echo/v4"
)

type OnlineStudentMessage struct {
	Action string `json:"action"`
	ID     string    `json:"user_id"`
}

type OnlineStudentsHandler struct {
	server *s.Server
}

func NewOnlineStudentsHandler(server *s.Server) *OnlineStudentsHandler {
	return &OnlineStudentsHandler{server: server}
}

func (handler *OnlineStudentsHandler) ConsumeOnlineStudentOld(c echo.Context) error {
	subscriber := redis_client.Init(handler.server.Config)
	userRepository := repositories.NewUserRepository(handler.server.DB)
	defer subscriber.Close()

	groupID := c.Param("group_id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	var userID []string
	userRepository.GetOnlineStudentsOld(&userID, groupID)

	fmt.Println(groupID)

	initialData, err := json.Marshal(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to marshal initial data")
	}
	fmt.Fprintf(c.Response().Writer, "data: %s\n\n", initialData)
	flusher.Flush()

	pubsub := subscriber.Subscribe(context.Background(), fmt.Sprintf("online-students:%s", groupID))
	defer pubsub.Close()

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var data OnlineStudentMessage
			if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
				fmt.Println(msg.Payload)
				fmt.Println("Failed to unmarshal message:", err)
				continue
			}

			switch data.Action {
			case "login":
				if !contains(userID, data.ID) {
					userID = append(userID, data.ID)
				}
			case "logout":
				userID = filter(userID, func(id string) bool { return id != data.ID })
			case "logout-all":
				userID = []string{}
			}

			updatedData, err := json.Marshal(userID)
			if err != nil {
				fmt.Println("Failed to marshal updated data:", err)
				continue
			}
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", updatedData)
			flusher.Flush()
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				userRepository.GetOnlineStudentsOld(&userID, groupID)
				tickerData, err := json.Marshal(userID)
				if err != nil {
					fmt.Println("Failed to marshal ticker data:", err)
					continue
				}
				fmt.Fprintf(c.Response().Writer, "data: %s\n\n", tickerData)
				flusher.Flush()
			case <-c.Request().Context().Done():
				return
			}
		}
	}()

	<- c.Request().Context().Done()

	pubsub.Unsubscribe(context.Background())
	pubsub.Close()
	
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}