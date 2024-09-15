package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   Server
	Postgres Postgres
}

type Server struct {
	ServerAddr string `env:"SERVER_ADDRESS" default:"0.0.0.0:8080"`
}

type Postgres struct {
	DSN string `env:"POSTGRES_CONN" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file reading error:", err.Error())
	}
	var cfg Config
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		panic("error while reading config: " + err.Error())
	}
	return &cfg
}
