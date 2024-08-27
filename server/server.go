package server

import (	
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-realtime-go/config"
	"github.com/Project-IPCA/ipca-realtime-go/db"
)

type Server struct {
	Echo    *echo.Echo
	DB      *gorm.DB
	Config  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.Init(cfg),
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}