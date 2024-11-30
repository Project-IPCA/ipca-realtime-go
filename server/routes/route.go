package routes

import (
	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/Project-IPCA/ipca-realtime-go/server/handlers"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ConfigureRoutes(server *s.Server){
	chapterPermissionHandler := handlers.NewChapterPermissionHandler(server)
	classLogHandler := handlers.NewClassLogHandler(server)
	onlineStudentsOldHandler := handlers.NewOnlineStudentsHandler(server)
	submitionResultHandler := handlers.NewSubmitionResultHandler(server)
	testCaseResultHandler := handlers.NewTestCaseResultHandler(server)
	userConnectionHandler := handlers.NewuserConnectionHandler(server)


	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())
	server.Echo.Use(middleware.CORS())

	apiGroup := server.Echo.Group("/subscribe")

	apiGroup.GET("/chapter-permission/:group_id",chapterPermissionHandler.ConsumeChapterPermission)
	apiGroup.GET("/class-log/:group_id",classLogHandler.ConsumeClassLog)
	apiGroup.GET("/online-students/:group_id",onlineStudentsOldHandler.ConsumeOnlineStudentOld)
	apiGroup.GET("/submission-result/:job_id",submitionResultHandler.ConsumeSubmitionResult)
	apiGroup.GET("/testcase-result/:job_id",testCaseResultHandler.ConsumeTestCaseResult)
	apiGroup.GET("/user/connection/:group_id",userConnectionHandler.ConsumeUserConnection)
}