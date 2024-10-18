package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"faraway/internal/config"
	"faraway/internal/domain"
	"fmt"
	"math"
)

// TODO May be split on 2 interfaces for Client (Solve) & Server (Challenge + Verify).
type POWServiceI interface {
	Challenge() []byte
	Verify(challenge, solution []byte) error
	Solve(token []byte) []byte
}

// POWService represents a proof of work algorithm implementation based on hashcash
type POWService struct {
	cfg config.POW
}

func NewPOWService(cfg config.POW) *POWService {
	return &POWService{cfg: cfg}
}

// Challenge returns a new token for a client.
func (p *POWService) Challenge() []byte {
	buf := make([]byte, p.cfg.TokenSize)

	const bits = 64

	// TODO need to add check on negative shift amount
	target := uint64(1) << (bits - p.cfg.Complexity)

	binary.BigEndian.PutUint64(buf[:8], target)
	_, _ = rand.Read(buf[8:])

	return buf
}

// Verify a client solution by its challenge.
func (p *POWService) Verify(challenge, solution []byte) error {
	if !verify(challenge, solution) {
		return fmt.Errorf("invalid solution: %w", domain.ErrInvalid)
	}

	return nil
}

// Solve a challenge & verify the solution.
func (p *POWService) Solve(token []byte) []byte {
	if len(token) != p.cfg.TokenSize {
		return nil
	}

	nonce := make([]byte, p.cfg.NonceSize)

	for i := uint64(0); i < math.MaxUint64; i++ {
		binary.BigEndian.PutUint64(nonce, i)

		if verify(token, nonce) {
			return nonce
		}
	}

	return nil
}

func verify(token, nonce []byte) bool {
	h := sha256.New()
	h.Write(token)
	h.Write(nonce)

	return bytes.Compare(h.Sum(nil), token) < 0
}
