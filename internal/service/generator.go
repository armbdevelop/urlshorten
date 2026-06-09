package service

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func generateShort() (string, error) {
	b := make([]byte, 10)
	maxIdx := big.NewInt(int64(len(alphabet)))

	for i := range b {
		index, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[index.Int64()]
	}
	return string(b), nil
}
