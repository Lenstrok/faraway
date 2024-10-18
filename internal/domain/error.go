package domain

import "errors"

var (
	ErrInvalid = errors.New("invalid")
)

type Error struct {
	Msg string `json:"message"`
}
