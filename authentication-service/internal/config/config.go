package config

import (
	"authentication-service/internal/models"
	"authentication-service/internal/repository"
)

type Config struct {
	DB     repository.DatabaseRepo
	Models models.Models
}
