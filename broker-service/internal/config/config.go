package config

import (
	"broker-service/event"
)

// Config is a config with all services that this app uses
type Config struct {
	AmqpService event.Service
}
