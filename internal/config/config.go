package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr        string `env:"ADDR" env-default:"localhost:8080"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	SecretKey   string `env:"SECRET_KEY" env-default:"very-secret-key"`
}

func MustLoadConfig() *Config {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		err := cleanenv.ReadConfig(".env", &cfg)
		if err != nil {
			log.Fatalf("Error reading .env: %v", err)
		}
	} else {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("Error reading: %v", err)
		}
	}

	return &cfg
}
