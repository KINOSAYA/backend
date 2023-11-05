package config

import (
	"authentication-service/internal/repository"
)

type Config struct {
	DB repository.DatabaseRepo
}
