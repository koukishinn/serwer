package security

import (
	"bytes"
	"encoding/csv"
	"os"
)

type SecureStore struct {
	file    string
	enclave *Enclave
	users   map[string]string
}

type SecureStoreOpts struct {
	File string
}

func NewSecureStore(opts *SecureStoreOpts) (*SecureStore, error) {
	data, err := os.ReadFile(opts.File)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(data))

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	users := make(map[string]string)

	for _, line := range lines {
		user := line[0]
		password := line[1]

		users[user] = password
	}

	return &SecureStore{
		file:    opts.File,
		users:   users,
		enclave: NewEnclave(),
	}, nil
}

func (s *SecureStore) Check(user, password string) (bool, error) {
	hash, err := s.enclave.Hash([]byte(password))
	if err != nil {
		return false, err
	}

	if uh, ok := s.users[user]; !ok || uh != hash {
		return false, nil
	}

	return true, nil
}
