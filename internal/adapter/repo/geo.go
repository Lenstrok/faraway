package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"vio_coding_challenge/internal/domain"

	"github.com/jmoiron/sqlx"
)

type GeoRepo struct {
	db *sqlx.DB
}

func NewGeoRepo(db *sqlx.DB) *GeoRepo {
	return &GeoRepo{db: db}
}

// SaveAll records in DB.
// On conflict by ip_address do nothing.
// Return the number of affected records.
// TODO It's good to add unittests with real DB or https://github.com/DATA-DOG/go-sqlmock.
func (gr GeoRepo) SaveAll(ctx context.Context, geos []domain.Geolocation) (int64, error) {
	query := `INSERT INTO geolocation (ip_address, country_code, country, city, latitude, longitude, mystery_value)
        	  VALUES (:ip_address, :country_code, :country, :city, :latitude, :longitude, :mystery_value)
        	  ON CONFLICT (ip_address) DO NOTHING`

	res, err := gr.db.NamedExecContext(ctx, query, geos)
	if err != nil {
		return 0, fmt.Errorf("failed to exec query: %w", err)
	}

	rowsCnt, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return rowsCnt, nil
}

// Get geolocation by ip address.
// Return domain.ErrNotFound in case of no such ip in database.
// TODO It's good to add unittests with real DB or https://github.com/DATA-DOG/go-sqlmock.
func (gr GeoRepo) Get(ctx context.Context, ipAddr string) (*domain.Geolocation, error) {
	geo := &domain.Geolocation{}

	query := `SELECT * FROM geolocation WHERE ip_address = $1`

	err := gr.db.GetContext(ctx, geo, query, ipAddr)
	if errors.Is(err, sql.ErrNoRows) {
		return geo, domain.ErrNotFound
	}

	return geo, err
}
