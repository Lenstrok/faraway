package main

import (
	"context"
	"faraway/internal/adapter/handler"
	"faraway/internal/adapter/repo"
	"faraway/internal/config"
	"faraway/internal/infra"
	"faraway/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	cfg := config.ServerConfig{}

	if err := config.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	quotesBytes, err := os.ReadFile(cfg.Quotes.FilePath)
	if err != nil {
		log.Fatalf("Failed to read file with quotes: %v", err)
	}

	quoteRepo, err := repo.NewQuoteRepo(quotesBytes)
	if err != nil {
		log.Fatalf("Failed to make quote repo: %v", err)
	}

	quoteService := service.NewQuoteService(quoteRepo)
	powService := service.NewPOWService(cfg.POW)
	netHandler := handler.NewServerNet(powService, quoteService)

	listener, err := infra.NewListener(netHandler, cfg.Server)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	log.Printf("Server started.")

	go listener.Start(ctx)

	// Shutdown waiting signal
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	log.Printf("Shutdown signal '%v' received.", <-interruptCh)
	close(interruptCh)

	listener.Stop()

	log.Printf("Server stopped.")
}
