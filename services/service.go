package services

import (
	"github.com/nexmedis-be-technical-test/repositories"
)

type Service struct {
	Repository *repositories.Repository
}

// NewService is the constructor for Service
func NewService(repository *repositories.Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
