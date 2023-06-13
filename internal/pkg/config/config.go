package config

import (
	"log"

	"github.com/caarlos0/env/v8"
)

var Config *Configuration

type Configuration struct {
	Server ServerConfiguration
}

type ServerConfiguration struct {
	Port     string `env:"PORT" envDefault:"8000"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func Setup() {
	var cfg *Configuration
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("couldn't parse configuration: %v", err)
	}

	Config = cfg
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
