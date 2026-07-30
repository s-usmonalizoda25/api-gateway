package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPORT     string `env:"HTTP_PORT" envDefault:"8080"`
	JWTSecretKey string `env:"JWT_SECRET_KEY" env-required:"true"`
	Services     Services
}

type Services struct {
	UserService    UserService    `env-prefix:"USER_SERVICE_"`
	MovieService   MovieService   `env-prefix:"MOVIE_SERVICE_"`
	BookingService BookingService `env-prefix:"BOOKING_SERVICE_"`
}

type UserService struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type MovieService struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type BookingService struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func New(path string) (*Config, error) {
	var conf Config

	err := cleanenv.ReadConfig(path, &conf)
	if err != nil {
		log.Printf("Config file at %s not found or error reading it, proceeding to read ENV", path)
	}

	err = cleanenv.ReadEnv(&conf)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return &conf, nil
}
