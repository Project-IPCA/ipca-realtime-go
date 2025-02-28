package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-realtime-go/models"
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

type UserConnection struct {
	UserID   string
	TabCount int
}

var (
	userConnections = make(map[string]*UserConnection)
	mu              sync.Mutex
)

func NewuserConnectionHandler(server *s.Server) *userConnectionHandler {
	return &userConnectionHandler{server: server}
}

type RedisMessage struct {
	Action string `json:"action"`
	UserID string `json:"user_id"`
}

func (userConnectionHanuserConnectionHandler *userConnectionHandler) ConsumeUserConnection(c echo.Context) error {
	pubsub := redis_client.Init(userConnectionHanuserConnectionHandler.server.Config)
	userRepository := repositories.NewUserRepository(userConnectionHanuserConnectionHandler.server.DB)
	defer pubsub.Close()

	userID := c.Param("user_id")

	var existsUser models.User
	userRepository.GetUser(&existsUser, userID)

	isStudent := false

	if existsUser.Student != nil {
		isStudent = true
	}

	mu.Lock()
	conn, exists := userConnections[userID]
	if !exists {
		conn = &UserConnection{
			UserID:   userID,
			TabCount: 1,
		}
		userConnections[userID] = conn
		userRepository.UpdateUserStatus(userID, true)
		if isStudent {
			message := RedisMessage{
				Action: "login",
				UserID: userID,
			}
			body, err := json.Marshal(message)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal message to JSON: "+err.Error())
			}
			err = pubsub.Publish(context.Background(), fmt.Sprintf("online-students:%s", existsUser.Student.GroupID), body).Err()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "can't publish message: "+err.Error())
			}
		}
	} else {
		conn.TabCount++
	}
	mu.Unlock()

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	userRepository.UpdateUserStatus(userID, true)

	subscriber := pubsub.Subscribe(context.Background(), fmt.Sprintf("user-event:%s", userID))
	defer subscriber.Close()

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported!")
	}

	go func() {
		ch := subscriber.Channel()
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

	mu.Lock()
	conn.TabCount--
	if conn.TabCount <= 0 {
		delete(userConnections, userID)
		userRepository.UpdateUserStatus(userID, false)
		fmt.Println("disconnect all")
		if isStudent {
			message := RedisMessage{
				Action: "logout",
				UserID: userID,
			}
			body, err := json.Marshal(message)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to marshal message to JSON: "+err.Error())
			}
			err = pubsub.Publish(context.Background(), fmt.Sprintf("online-students:%s", existsUser.Student.GroupID), body).Err()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "can't publish message: "+err.Error())
			}
		}
	}
	mu.Unlock()

	subscriber.Unsubscribe(context.Background())
	subscriber.Close()

	return nil
}
