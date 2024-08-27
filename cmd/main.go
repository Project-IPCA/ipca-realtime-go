package main

import (
	"fmt"

	application "github.com/Project-IPCA/ipca-realtime-go"
	"github.com/Project-IPCA/ipca-realtime-go/config"
	"github.com/Project-IPCA/ipca-realtime-go/docs"
)

// @title IPCA API
// @version 1.0
// @description IPCA API
// @host petstore.swagger.io
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	cfg := config.NewConfig()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	application.Start(cfg)
}
