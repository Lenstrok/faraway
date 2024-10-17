package main

import (
	"context"
	"log"

	"vio_coding_challenge/internal/adapter/repo"
	"vio_coding_challenge/internal/config"
	"vio_coding_challenge/internal/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dbDriverName = "postgres"

func main() {
	ctx := context.Background()

	cfg := config.ParserConfig{}

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

	log.Printf("CSV export started.")

	stats, err := geoService.ParseCSV(ctx, cfg.Parser.FilePath, cfg.Parser.BatchSize)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	log.Printf(
		"CSV export completed! Stats: accepted_rows=%d, discarded_rows=%d, time_elapsed=%s",
		stats.RecAcceptedCnt, stats.RecDiscardedCnt, stats.TimeElapsed,
	)
}
