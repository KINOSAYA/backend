package config

import (
	"broker-service/internal/routers"
)

func NewConfig(router routers.Router) *Config {
	return &Config{
		Router: router,
	}
}

type Config struct {
	Router routers.Router
}
