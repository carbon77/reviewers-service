package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port       int    `env:"PORT"`
	DbUser     string `env:"DB_USER,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbName     string `env:"DB_NAME,required"`
	DbHost     string `env:"DB_HOST"`
	DbPort     int    `env:"DB_PORT"`
}

func Load() *Config {
	cfg := Config{
		Port:   8080,
		DbHost: "localhost",
		DbPort: 5432,
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Can't parse configuration")
	}

	return &cfg
}
