package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type MigrationConfig struct {
	DB DB
}

type ServerConfig struct {
	DB     DB
	Server Server
}

type ParserConfig struct {
	DB     DB
	Parser Parser
}

func NewConfig[T ServerConfig | ParserConfig | MigrationConfig](cfg *T) error {
	if err := godotenv.Load("env_example"); err != nil {
		log.Printf("Try to get env without .env file: %v", err)
	}

	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("failed to procees config: %w", err)
	}

	return nil
}
