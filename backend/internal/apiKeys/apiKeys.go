package apikeys

import (
	"errors"

	"github.com/relextm19/tracker.nvim/internal/helpers"
	"github.com/relextm19/tracker.nvim/internal/users"
)

type APIKey struct {
	Name    string `json:"name"`
	KeyHash string `json:"keyHash"`
}

type ClientAPIKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func NewClientAPIKey() *ClientAPIKey {
	return &ClientAPIKey{}
}

func (cak *ClientAPIKey) Valid() error {
	if !helpers.ValidStringField(cak.Name) {
		return errors.New("api key name is required")
	}
	if !helpers.ValidStringField(cak.Key) {
		return errors.New("api key is required")
	}

	return nil
}

func NewAPIKey(cak *ClientAPIKey) (*APIKey, error) {
	keyHash, err := users.GetHash(cak.Key)
	if err != nil {
		return nil, err
	}

	return &APIKey{
		Name:    cak.Name,
		KeyHash: keyHash,
	}, nil
}
