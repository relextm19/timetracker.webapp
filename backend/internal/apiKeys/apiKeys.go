package apikeys

import (
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
