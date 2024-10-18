package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ClientConfig struct {
	POW    POW
	Server Server
}

type ServerConfig struct {
	POW    POW
	Quotes Quotes
	Server Server
}

func NewConfig[T ServerConfig | ClientConfig](cfg *T) error {
	if err := godotenv.Load("env.example"); err != nil {
		log.Printf("Try to get env without .env file: %v", err)
	}

	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("failed to procees config: %w", err)
	}

	return nil
}
