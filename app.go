package application

import (
	"log"

	"github.com/Project-IPCA/ipca-realtime-go/config"
	"github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/Project-IPCA/ipca-realtime-go/server/routes"
)

func Start(cfg *config.Config) {
	app := server.NewServer(cfg)

	routes.ConfigureRoutes(app)

	err := app.Start(cfg.HTTP.Port)
	if err != nil {
		log.Fatal("Port already used")
	}
}
