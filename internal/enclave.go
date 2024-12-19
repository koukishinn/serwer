package internal

import (
	"crypto/rand"
	"encoding/hex"
)

type Enclave struct{}

func (e *Enclave) Generate() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
