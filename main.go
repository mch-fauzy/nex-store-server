package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/nexmedis-be-technical-test/configs"
	"github.com/nexmedis-be-technical-test/handlers"
	"github.com/nexmedis-be-technical-test/repositories"
	"github.com/nexmedis-be-technical-test/routes"
	"github.com/nexmedis-be-technical-test/services"
)

func main() {
	// Initialize logger
	configs.InitLogger()

	// Initialize the PostgreSQL connection
	config := configs.Get()
	dbConn := configs.NewPostgreSqlConn(config)

	// Initialize repository, service, and handler layers
	repository := repositories.NewRepository(dbConn)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service)

	// Setup the router
	route := routes.SetupRouter(handler)

	// Start the server
	log.Info().Str("port", config.Server.Port).Msg("Starting up HTTP server")
	err := http.ListenAndServe(":"+config.Server.Port, route)
	if err != nil {
		log.Error().Err(err).Msg("Server failed to start")
	}
}
