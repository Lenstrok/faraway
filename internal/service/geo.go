package service

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"vio_coding_challenge/internal/domain"
	"vio_coding_challenge/internal/service/port"

	"github.com/go-playground/validator/v10"
)

// ParsingStats for CSV about the time elapsed,
// as well as the number of entries accepted/discarded.
type ParsingStats struct {
	RecAcceptedCnt  int
	RecDiscardedCnt int

	TimeElapsed time.Duration
}

//go:generate mockgen -destination=../../mocks/service/mock_geo.go -package=mocks -source=geo.go GeoServiceI

type GeoServiceI interface {
	GetByIP(ctx context.Context, ip string) (*domain.Geolocation, error)
	ParseCSV(ctx context.Context, filePath string, batchSize int) (*ParsingStats, error)
}

// GeoService provides import and export of geolocation data.
type GeoService struct {
	geoRepo port.GeoRepoI

	validate *validator.Validate
}

func NewGeoService(geoRepo port.GeoRepoI, validate *validator.Validate) *GeoService {
	return &GeoService{geoRepo: geoRepo, validate: validate}
}

// GetByIP a corresponding geolocation record.
func (gs GeoService) GetByIP(ctx context.Context, ip string) (*domain.Geolocation, error) {
	return gs.geoRepo.Get(ctx, ip)
}

// ParseCSV containing the raw data and persists it in a database.
// filePath param is absolute CSV file path.
//
// Sanitise the entries by mapping with domain.Geolocation and validating the struct leveraging `validate` tag.
// All date constraints described in the domain.Geolocation struct.
// IP Address duplication leveraging by DB (Primary Key).
//
// Store data in DB by batches. batchSize param define exact size of such batch.
//
// At the end of the parsing process, return statistics about the time elapsed,
// as well as the number of entries accepted/discarded.
func (gs GeoService) ParseCSV(ctx context.Context, filePath string, batchSize int) (*ParsingStats, error) {
	startTime := time.Now()
	stats := &ParsingStats{}

	file, err := os.Open(filePath)
	if err != nil {
		return stats, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			log.Printf("Failed to close csv file: %v", cErr)
		}
	}()

	reader := csv.NewReader(bufio.NewReader(file))

	headers, err := reader.Read()
	if err != nil {
		return stats, fmt.Errorf("failed to read headers: %w", err)
	}

	const columnsCnt = 7

	if len(headers) != columnsCnt {
		return stats, fmt.Errorf("expected %d columns, got %d", columnsCnt, len(headers))
	}

	geoBatch := make([]domain.Geolocation, 0)

	for {
		record, rErr := reader.Read()
		if rErr != nil {
			if rErr == io.EOF {
				break
			}

			return stats, fmt.Errorf("failed to read csv record: %w", rErr)
		}

		geo, nErr := gs.parseRecord(record)
		if nErr != nil {
			stats.RecDiscardedCnt++
			continue
		}

		geoBatch = append(geoBatch, geo)

		if len(geoBatch) >= batchSize {
			if err = gs.saveBatch(ctx, geoBatch, stats); err != nil {
				return stats, fmt.Errorf("failed to save batch: %w", err)
			}

			geoBatch = geoBatch[:0]
		}
	}

	if len(geoBatch) > 0 {
		if err = gs.saveBatch(ctx, geoBatch, stats); err != nil {
			return stats, fmt.Errorf("failed to save batch: %w", err)
		}
	}

	stats.TimeElapsed = time.Since(startTime)

	return stats, nil
}

// SaveBatch of records in the DB.
// All accepted records add to the RecAcceptedCnt stats.
// On conflict with ip_address PK do nothing & add these records to the RecDiscardedCnt stats.
func (gs GeoService) saveBatch(ctx context.Context, batch []domain.Geolocation, stats *ParsingStats) error {
	rowsCnt, err := gs.geoRepo.SaveAll(ctx, batch)
	if err != nil {
		return err
	}

	stats.RecAcceptedCnt += int(rowsCnt)
	stats.RecDiscardedCnt += len(batch) - int(rowsCnt)

	return nil
}

// parseRecord of 7 elements and return domain.Geolocation.
// Return domain.ErrInvalid in case of validation errors.
func (gs GeoService) parseRecord(record []string) (domain.Geolocation, error) {
	// TODO We can add size restrictions on values in each column.
	geo := domain.Geolocation{
		IPAddress:   record[0],
		CountryCode: record[1],
		Country:     record[2],
		City:        record[3],
		// TODO don't have enough info about MysteryValue type & left it as a string (text in DB).
		MysteryValue: record[6],
	}

	var err error
	if geo.Latitude, err = strconv.ParseFloat(record[4], 64); err != nil {
		return domain.Geolocation{}, domain.ErrInvalid
	}

	if geo.Longitude, err = strconv.ParseFloat(record[5], 64); err != nil {
		return domain.Geolocation{}, domain.ErrInvalid
	}

	if err = gs.validate.Struct(&geo); err != nil {
		return domain.Geolocation{}, domain.ErrInvalid
	}

	return geo, nil
}
