package service

import (
	"crypto/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func generateShort() (string, error) {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// берем остаток от деления
	// У нас алфавит из 63 символов. Нам нужно любое число от 0 до 255 превратить в число от 0 до 62 — индекс в
	//  алфавите. Остаток от деления (%) как раз даёт число от 0 до 62.
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}

	return string(b), nil
}
