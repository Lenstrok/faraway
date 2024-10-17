package config

import "fmt"

const dbDSNFmt = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"

type DB struct {
	User string `envconfig:"DB_USER" default:"admin"`
	Pass string `envconfig:"DB_PASSWORD" default:"admin"`
	Name string `envconfig:"DB_NAME" default:"vio_challenge"`
	Host string `envconfig:"DB_HOST" default:"localhost"`
	Port string `envconfig:"DB_PORT" default:"5432"`
}

func (db DB) GetDSN() string {
	return fmt.Sprintf(dbDSNFmt, db.Host, db.User, db.Pass, db.Name, db.Port)
}
