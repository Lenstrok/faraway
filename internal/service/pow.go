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

// todo rename params + refactor + tests

type POWServiceI interface {
	// todo rename methods
	Challenge() []byte
	Verify(challenge, solution []byte) error
	Solve(token []byte) []byte
}

// todo update POW represents a proof of work algorithm implementation based on hashcash
type POWService struct {
	cfg config.POW
}

// TODO
func NewPOWService(cfg config.POW) *POWService {
	return &POWService{cfg: cfg}
}

// todo docs
func (p *POWService) Challenge() []byte {
	buf := make([]byte, p.cfg.TokenSize)
	target := uint64(1) << (64 - p.cfg.Complexity)

	binary.BigEndian.PutUint64(buf[:8], target)
	_, _ = rand.Read(buf[8:])

	return buf
}

// todo docs
func (p *POWService) Verify(challenge, solution []byte) error {
	if !verify(challenge, solution) {
		return fmt.Errorf("invalid solution: %w", domain.ErrInvalid)
	}

	return nil
}

// TODO docs + test
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

// todo docs
func verify(token, nonce []byte) bool {
	h := sha256.New()
	h.Write(token)
	h.Write(nonce)

	return bytes.Compare(h.Sum(nil), token) < 0
}
