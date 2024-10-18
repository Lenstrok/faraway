package main

import (
	"context"
	"faraway/internal/adapter/handler"
	"faraway/internal/config"
	"faraway/internal/infra"
	"faraway/internal/service"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	cfg := config.ClientConfig{}

	if err := config.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	powService := service.NewPOWService(cfg.POW)
	netHandler := handler.NewClientNet(powService)
	dialer := infra.NewDialer(netHandler, cfg.Server)

	log.Printf("Client started.")

	for range cfg.NumberOfQuotes {
		quote, err := dialer.GetQuote(ctx)
		if err != nil {
			log.Printf("Failed to handle quote request: %v", err)
		}

		log.Printf("Got a quote: %s", quote)
	}

	log.Printf("Client stopped.")
}
