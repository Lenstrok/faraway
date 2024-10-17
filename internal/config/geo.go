package config

type Parser struct {
	// BatchSize determines how many records we send to the database at the same time.
	BatchSize int `envconfig:"PARSER_BATCH_SIZE" default:"5000"`
	// FilePath is an absolute path to the file.
	FilePath string `envconfig:"PARSER_FILE_PATH" required:"true"`
}
