package apikeys

import (
	"errors"

	"github.com/google/uuid"
	"github.com/relextm19/tracker.nvim/internal/helpers"
)

type RequestAPIKey struct {
	Name string `json:"name"`
}

type APIKey struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int    `json:"createdAt"`
	KeyHash   string `json:"keyHash"`
	Key       string `json:"key"`
}

func NewRequestAPIKey() *RequestAPIKey {
	return &RequestAPIKey{}
}

func (rak *RequestAPIKey) Valid() error {
	if !helpers.ValidStringField(rak.Name) {
		return errors.New("api key name is required")
	}

	return nil
}

// NewAPIKey Only fills out values know at the time of the function call you have to fill out the rest after querying the db
func NewAPIKey(cak *RequestAPIKey) (*APIKey, error) {
	key := uuid.New().String()
	keyHash, err := helpers.GetHash([]byte(key))
	if err != nil {
		return nil, err
	}

	return &APIKey{
		Name:    cak.Name,
		KeyHash: keyHash,
		Key:     key,
	}, nil
}
