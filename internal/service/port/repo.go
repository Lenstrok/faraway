package port

import (
	"context"
	"vio_coding_challenge/internal/domain"
)

//go:generate mockgen -destination=../../../mocks/repo/mock_geo.go -package=mocks -source=repo.go GeoRepoI

type GeoRepoI interface {
	SaveAll(ctx context.Context, geos []domain.Geolocation) (int64, error)
	Get(ctx context.Context, ipAddr string) (*domain.Geolocation, error)
}
