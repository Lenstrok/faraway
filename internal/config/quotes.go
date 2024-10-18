package config

type Quotes struct {
	FilePath string `envconfig:"QUOTES_FILE_PATH" default:"internal/config/quotes.json"`
}
