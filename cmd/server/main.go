package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	v1 "vio_coding_challenge/internal/adapter/controller/v1"
	"vio_coding_challenge/internal/adapter/repo"
	"vio_coding_challenge/internal/config"
	"vio_coding_challenge/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// main
// Swagger spec:
// @title       Vio Challenge Server
// @description Here we handle Geolocation!
// @version     1.0
// @BasePath    /
func main() {
	ctx := context.Background()

	cfg := config.ServerConfig{}

	err := config.NewConfig(&cfg)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	// For production run we should add lifetime & conn pool size.
	db, err := sqlx.Connect(dbDriverName, cfg.DB.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	geoRepo := repo.NewGeoRepo(db)
	geoService := service.NewGeoService(geoRepo, validator.New(validator.WithRequiredStructEnabled()))
	geoCtrl := v1.NewGeoCtrl(geoService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1.AddRoutes(router, geoCtrl)

	server := &http.Server{
		Handler:           router,
		Addr:              ":" + cfg.Server.Port,
		ReadTimeout:       cfg.Server.Timeout,
		ReadHeaderTimeout: cfg.Server.Timeout,
		WriteTimeout:      cfg.Server.Timeout,
	}

	go func() {
		_ = server.ListenAndServe()
	}()

	log.Printf("Server started.")

	// Shutdown waiting signal
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	log.Printf("Shutdown signal '%v' received.", <-interruptCh)
	close(interruptCh)

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("Failed to stop server: %v", err)
	}

	log.Printf("Server stopped.")
}
