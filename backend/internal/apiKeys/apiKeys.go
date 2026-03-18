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

type RequestAPIKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ResponseAPIKey struct {
	Name      string `json:"name"`
	CreatedAt int    `json:"createdAt"`
}

func NewResponseAPIKey() *ResponseAPIKey {
	return &ResponseAPIKey{}
}

func NewRequestAPIKey() *RequestAPIKey {
	return &RequestAPIKey{}
}

func (rak *RequestAPIKey) Valid() error {
	if !helpers.ValidStringField(rak.Name) {
		return errors.New("api key name is required")
	}
	if !helpers.ValidStringField(rak.Key) {
		return errors.New("api key is required")
	}

	return nil
}

func NewAPIKey(cak *RequestAPIKey) (*APIKey, error) {
	keyHash, err := users.GetHash(cak.Key)
	if err != nil {
		return nil, err
	}

	return &APIKey{
		Name:    cak.Name,
		KeyHash: keyHash,
	}, nil
}
