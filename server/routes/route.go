package routes

import (
	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/Project-IPCA/ipca-realtime-go/server/handlers"
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

	apiGroup := server.Echo.Group("/subscribe")

	apiGroup.GET("/group-permission/:group_id", groupPermissionHandler.ConsumeGroupPermission)
	apiGroup.GET("/class-log/:group_id", classLogHandler.ConsumeClassLog)
	apiGroup.GET("/online-students/:group_id", onlineStudentsOldHandler.ConsumeOnlineStudentOld)
	apiGroup.GET("/submission-result/:job_id", submitionResultHandler.ConsumeSubmitionResult)
	apiGroup.GET("/testcase-result/:job_id", testCaseResultHandler.ConsumeTestCaseResult)
	apiGroup.GET("/user/connection/:user_id", userConnectionHandler.ConsumeUserConnection)
}
