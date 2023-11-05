package config

import (
	"broker-service/internal/routes"
)

func NewConfig(router routes.Router) *Config {
	return &Config{
		Router: router,
	}
}

type Config struct {
	Router routes.Router
}
