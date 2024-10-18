package main

import (
	"context"
	"faraway/internal/adapter/controller"
	"faraway/internal/adapter/repo"
	"faraway/internal/config"
	"faraway/internal/service"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net"
	"os"
)

func main() {
	ctx := context.Background()

	cfg := config.ServerConfig{}

	if err := config.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	quotesBytes, err := os.ReadFile("internal/config/quotes.json") // todo to conf?
	if err != nil {
		log.Fatalf("Failed to read file with quotes: %v", err)
	}

	quoteRepo, err := repo.NewQuoteRepo(quotesBytes)
	if err != nil {
		log.Fatalf("Failed to make quote repo: %v", err)
	}

	quoteService := service.NewQuoteService(quoteRepo)
	powService, err := service.NewPOWService(10) // todo to cfg?
	if err != nil {
		log.Fatalf("Failed to make pow service: %v", err)
	}

	handler := controller.NewHandler(powService, quoteService)

	// todo describe about moving to the infra level
	listener, err := net.Listen("tcp", ":8080") // todo const/conf
	if err != nil {
		log.Fatalf("net listen error: %v\n", err)
	}

	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			log.Printf("Failed to close listener: %v", closeErr)
		}
	}()

	// todo добавить асинхрон или упомянуть
	// todo добавить шатдаун
	for {
		conn, AcceptErr := listener.Accept()
		if AcceptErr != nil {
			log.Fatalf("accept connection error: %v\n", AcceptErr)
		}

		go func() {
			if handleErr := handler.HandleQuoteReq(ctx, conn); handleErr != nil {
				log.Printf("Failed to handle quote request: %v", handleErr)
			}
		}()
	}

	// todo add more logs
	log.Printf("Server stopped.")
}
