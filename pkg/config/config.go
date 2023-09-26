package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL"`
	Server   struct {
		Port int    `env:"PORT"`
		Host string `env:"HOST"`
	}
	RabbitMQ struct {
		Host     string `env:"RMQ_HOST"`
		Port     string `env:"RMQ_PORT"`
		Username string `env:"RMQ_USER"`
		Password string `env:"RMQ_PASSWORD"`
	}
	N int `env:"CONSUMER_NUMBER" envDefault:"1"`
}

func Init() *Config {
	config := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	if err = env.Parse(&config); err != nil {
		log.Fatal(err)
	}

	return &config
}
