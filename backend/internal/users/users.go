package users

import (
	"errors"
	"net/mail"

	"github.com/google/uuid"
	"github.com/relextm19/tracker.nvim/internal/helpers"
)

type RequestUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRequestUserBody() *RequestUserBody {
	return &RequestUserBody{}
}

type User struct {
	Email        string
	PasswordHash string
	Token        uuid.UUID
}

func NewUser(cub *RequestUserBody) (*User, error) {
	ph, err := helpers.GetHashFromString(cub.Password)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:        cub.Email,
		PasswordHash: ph,
		Token:        uuid.New(),
	}, nil
}

func (cub *RequestUserBody) Valid() error {
	if ok := helpers.ValidStringField(cub.Email); !ok {
		return helpers.ErrEmptyField
	}

	if _, err := mail.ParseAddress(cub.Email); err != nil {
		return errors.New("invalid email format")
	}

	if ok := helpers.ValidStringField(cub.Password); !ok {
		return helpers.ErrEmptyField
	}

	return nil
}
