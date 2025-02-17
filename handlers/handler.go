package handlers

import (
	"github.com/nexmedis-be-technical-test/services"
)

type Handler struct {
	Service *services.Service
}

// NewHandler is the constructor for Handler
func NewHandler(service *services.Service) Handler {
	return Handler{
		Service: service,
	}
}
