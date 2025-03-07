package routes

import (
	"github.com/Project-IPCA/ipca-realtime-go/middlewares"
	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/Project-IPCA/ipca-realtime-go/server/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server) {
	groupPermissionHandler := handlers.NewGroupPermissionHandler(server)
	classLogHandler := handlers.NewClassLogHandler(server)
	onlineStudentsOldHandler := handlers.NewOnlineStudentsHandler(server)
	submitionResultHandler := handlers.NewSubmitionResultHandler(server)
	testCaseResultHandler := handlers.NewTestCaseResultHandler(server)
	userConnectionHandler := handlers.NewuserConnectionHandler(server)

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())
	server.Echo.Use(middleware.CORS())

	authMiddleware := middlewares.NewAuthMiddleware(server)
	jwtConfig := authMiddleware.GetJwtConfig()

	apiGroup := server.Echo.Group("/subscribe")
	apiAuthGroup := apiGroup
	apiAuthGroup.Use(echojwt.WithConfig(jwtConfig))

	apiAuthGroup.GET("/group-permission/:group_id", groupPermissionHandler.ConsumeGroupPermission)
	apiAuthGroup.GET("/class-log/:group_id", classLogHandler.ConsumeClassLog)
	apiAuthGroup.GET("/online-students/:group_id", onlineStudentsOldHandler.ConsumeOnlineStudent)
	apiAuthGroup.GET("/submission-result/:job_id", submitionResultHandler.ConsumeSubmitionResult)
	apiAuthGroup.GET("/testcase-result/:job_id", testCaseResultHandler.ConsumeTestCaseResult)
	apiAuthGroup.GET("/user/connection/:user_id", userConnectionHandler.ConsumeUserConnection)
}
