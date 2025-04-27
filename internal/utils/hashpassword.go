package utils

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// sturct to store the hash and salt
type argon2idHash struct {
	params *argon2id.Params
}

// define an interface for the hashed password
type HashedPassword interface {
	Compare(password string, storedhash string) (bool, error)
	Hash(password string) (string, error)
}

// instance of the hashed password
func NewHashedPassword(params *argon2id.Params) HashedPassword {
	if params == nil {
		params = argon2id.DefaultParams
	}
	return &argon2idHash{params: params}
}

func (h *argon2idHash) Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, h.params)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return hash, nil
}

func (h *argon2idHash) Compare(password string, storedhash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, storedhash)
	if err != nil {
		return false, fmt.Errorf("invalid password: %w", err)
	}
	return match, nil

}
