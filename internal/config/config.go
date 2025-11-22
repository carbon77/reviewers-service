package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HttpPort string `env:"PORT"`
}

func Load() *Config {
	cfg := Config{
		HttpPort: "8080",
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Can't parse configuration")
	}

	return &cfg
}
