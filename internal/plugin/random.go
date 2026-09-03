package plugin

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type randomSource interface {
	Intn(max int) (int, error)
}

type cryptoRandom struct{}

func (cryptoRandom) Intn(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("random upper bound must be positive")
	}
	value, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("read cryptographic random source: %w", err)
	}
	return int(value.Int64()), nil
}
