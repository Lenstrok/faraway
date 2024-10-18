package main

import (
	"context"
	"faraway/internal/adapter/controller"
	"faraway/internal/config"
	"faraway/internal/service"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	cfg := config.ServerConfig{}

	if err := config.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	powService, err := service.NewPOWService(10) // todo to cfg?
	if err != nil {
		log.Fatalf("Failed to make pow service: %v", err)
	}

	dialerHandler := controller.NewDialer(powService) // todo rename package

	// todo describe about moving to the infra level
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", ":8080") // todo const/conf
	if err != nil {
		log.Fatalf("Failed to dial conn: %v", err)
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("Failed to close listener: %v", closeErr)
		}
	}()

	quote, err := dialerHandler.GetQuote(ctx, conn)
	if err != nil {
		log.Printf("Failed to handle quote request: %v", err)
	}

	log.Printf("Got a quote: %s", quote) // todo

	// todo добавить асинхрон или упомянуть
	// todo добавить шатдаун

	// todo add more logs
	log.Printf("Server stopped.")
}
