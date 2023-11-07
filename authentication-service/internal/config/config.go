package config

import (
	"authentication-service/internal/models"
	"authentication-service/internal/repository"
	"authentication-service/internal/service"
)

type Config struct {
	DB      repository.DatabaseRepo
	Models  models.Models
	Service service.AuthorizationService
}
