package domain

type Geolocation struct {
	IPAddress   string  `json:"ip_address" db:"ip_address" validate:"required,ip_addr"`
	CountryCode string  `json:"country_code" db:"country_code" validate:"required,country_code"`
	Country     string  `json:"country" db:"country" validate:"required,lte=255"`
	City        string  `json:"city" db:"city" validate:"required,lte=255"`
	Latitude    float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" db:"longitude" validate:"required"`
	// TODO don't have enough info about MysteryValue type & left it as a string (text in DB).
	MysteryValue string `json:"mystery_value" db:"mystery_value" validate:"required"`
}
