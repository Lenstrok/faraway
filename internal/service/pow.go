package service

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"faraway/internal/domain"
	"fmt"
	"math"
)

// todo rename params + refactor + tests

// TODO change + to cfg
const (
	tokenSize = 16
	nonceSize = 8
)

type POWServiceI interface {
	// todo rename methods
	Challenge() []byte
	Verify(challenge, solution []byte) error
	Solve(token []byte) []byte
}

// todo update POW represents a proof of work algorithm implementation based on hashcash
type POW struct {
	complexity uint64
}

// TODO
func NewPOW(complexity uint64) (*POW, error) {
	const maxTargetBits = 24 // todo cfg?

	if complexity < 1 || complexity > maxTargetBits {
		return nil, fmt.Errorf("invalid complexity value: %w", domain.ErrInvalid)
	}

	return &POW{complexity: complexity}, nil
}

// todo docs
func (p *POW) Challenge() []byte {
	buf := make([]byte, tokenSize)
	target := uint64(1) << (64 - p.complexity)

	binary.BigEndian.PutUint64(buf[:8], target)
	_, _ = rand.Read(buf[8:])

	return buf
}

// todo docs
func (p *POW) Verify(challenge, solution []byte) error {
	if len(challenge) != tokenSize {
		return fmt.Errorf("invalid challenge size: %w", domain.ErrInvalid)
	}

	if len(solution) != nonceSize {
		return fmt.Errorf("invalid solution size: %w", domain.ErrInvalid)
	}

	if !verify(challenge, solution) {
		return fmt.Errorf("invalid solution: %w", domain.ErrInvalid)
	}

	return nil
}

// TODO docs + test
func (p *POW) Solve(token []byte) []byte {
	if len(token) != tokenSize {
		return nil
	}

	nonce := make([]byte, nonceSize)

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
