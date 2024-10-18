package config

import "time"

type Server struct {
	Port        string        `envconfig:"SERVER_PORT" default:"8080"`
	ConnTimeout time.Duration `envconfig:"CONN_TIMEOUT" default:"10s"`
}
