package security

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

type Enclave struct {
	hasher hash.Hash
}

func NewEnclave() *Enclave {
	hasher := sha512.New()

	return &Enclave{
		hasher: hasher,
	}
}

func (e *Enclave) Nonce() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (e *Enclave) Hash(s []byte) (string, error) {
	_, err := e.hasher.Write(s)
	if err != nil {
		return "", err
	}

	hashed := hex.EncodeToString(e.hasher.Sum(nil))

	e.hasher.Reset()

	return hashed, nil
}
