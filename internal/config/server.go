package config

import "time"

type Server struct {
	Port    string        `envconfig:"SERVER_PORT" default:"8080"`
	Timeout time.Duration `envconfig:"SERVER_TIMEOUT" default:"5s"`
}
