package config

type POW struct {
	TokenSize  int `envconfig:"POW_TOKEN_SIZE" default:"16"`
	NonceSize  int `envconfig:"POW_NONCE_SIZE" default:"8"`
	Complexity int `envconfig:"POW_COMPLEXITY" default:"10"`
}
